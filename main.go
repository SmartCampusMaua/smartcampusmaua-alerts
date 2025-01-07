package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"crypto/tls"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	// "github.com/joho/godotenv"
)

type SmartLightFields struct {
	BatteryVoltage float64 `json:"batteryVoltage"`
	BoardVoltage   float64 `json:"boardVoltage"`
	Data           string  `json:"data"`
	FCnt           float64 `json:"fCnt"`
	FPort          float64 `json:"fPort"`
	Humidity       float64 `json:"humidity"`
	Luminosity     float64 `json:"luminosity"`
	Movement       float64 `json:"movement"`
	RxAlt_0        float64 `json:"rxAlt_0"`
	RxLat_0        float64 `json:"rxLat_0"`
	RxLon_0        float64 `json:"rxLon_0"`
	RxRssi_0       float64 `json:"rxRssi_0"`
	RxSnr_0        float64 `json:"rxSnr_0"`
	Temperature    float64 `json:"temperature"`
	TxBandWidth    float64 `json:"txBandWidth"`
	TxFrequency    float64 `json:"txFrequency"`
	TxSpreadFactor float64 `json:"txSpreadFactor"`
}

type SmartLightTags struct {
	DeviceId     string `json:"deviceId"`
	DeviceType   string `json:"deviceType"`
	Direction    string `json:"direction"`
	Host         string `json:"host"`
	Origin       string `json:"origin"`
	RxMac_0      string `json:"rxMac_0"`
	TxCodeRate   string `json:"txCodeRate"`
	TxModulation string `json:"txModulation"`
	Type         string `json:"type"`
}

type SmartLightData struct {
	Fields    SmartLightFields `json:"fields"`
	Name      string           `json:"name"`
	Tags      SmartLightTags   `json:"tags"`
	Timestamp float64          `json:"timestamp"`
}

func fetchSmartLight() ([]SmartLightData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/SmartLight/all?interval=300"

	// Create a custom HTTP client that doesn't verify SSL certificates
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Disable SSL verification
			},
		},
		Timeout: 30 * time.Second, // Optional timeout for the request
	}

	// Make the HTTP GET request
	response, err := client.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer response.Body.Close()

	// Check for successful HTTP response status
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response into a slice of SmartLightData (since the response is an array)
	var data []SmartLightData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]SmartLightData, 0)
	seenDevices := make(map[string]bool)

	for _, item := range data {
		if !seenDevices[item.Tags.DeviceId] {
			item.Fields.BatteryVoltage = item.Fields.BatteryVoltage / 1000 // fix temporario até arrumar o timeseries
			uniqueData = append(uniqueData, item)
			seenDevices[item.Tags.DeviceId] = true
		}
	}

	// Return the filtered data
	return uniqueData, nil
}

func fetchUsers() ([]UserData, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	connStr := os.Getenv("supabaseConnection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	query := `SELECT "alarms", "phone" FROM "User"`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying the database: %v", err)
	}
	defer rows.Close()

	var userData []UserData

	for rows.Next() {
		var loc UserData
		var alarmsTemp string

		err := rows.Scan(&alarmsTemp, &loc.Phone)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if alarmsTemp != "{}" {
			alarmsTemp = strings.ReplaceAll(alarmsTemp, "\\\"", "\"")
			alarmsTemp = "" + alarmsTemp[2:len(alarmsTemp)-2] + ""
			alarmsTemp = "[" + alarmsTemp + "]"
			alarmsTemp = strings.Replace(alarmsTemp, `","`, ",", -1)

			err := json.Unmarshal([]byte(alarmsTemp), &loc.Alarms)
			if err != nil {
				log.Fatal(err)
			}

			userData = append(userData, loc)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return userData, nil
}

func updateAlarmAlreadyPlayedOnSupabase(messages []Message) {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	connStr := os.Getenv("supabaseConnection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	updateQuery := `
		UPDATE "User"
		SET "alarms" = array(
			SELECT CASE
					WHEN element->>'deveui' = $1
						AND element->>'trigger' = $2
						AND element->>'triggerType' = $3
						AND element->>'triggerAt' = $4
					THEN jsonb_set(element, '{alreadyPlayed}', 'true'::jsonb)
					ELSE element
			END
			FROM unnest("alarms") AS element
		)
	`

	for _, msg := range messages {

		_, err := db.Exec(updateQuery, msg.DEVEUI, msg.Trigger, msg.TriggerType, msg.TriggerAt)
		if err != nil {
			log.Printf("Error updating alarm for DEVEUI: %s, Trigger: %s, TriggerType: %s, TriggerAt: %s. Error: %v",
				msg.DEVEUI, msg.Trigger, msg.TriggerType, msg.TriggerAt, err)
		} else {
			fmt.Printf("Successfully updated alarm for DEVEUI: %s, Trigger: %s, TriggerType: %s, TriggerAt: %s\n",
				msg.DEVEUI, msg.Trigger, msg.TriggerType, msg.TriggerAt)
		}
	}
}

type UserData struct {
	Alarms []Alarm `json:"alarms"`
	Phone  string  `json:"phone"`
}

type Alarm struct {
	Type          string `json:"type"`
	Local         string `json:"local"`
	DEVEUI        string `json:"deveui"`
	Trigger       string `json:"trigger"`
	TriggerType   string `json:"triggerType"`
	TriggerAt     string `json:"triggerAt"`
	AlreadyPlayed bool   `json:"alreadyPlayed"`
}

type Message struct {
	DEVEUI       string
	Type         string
	Trigger      string
	TriggerType  string
	TriggerAt    string
	Phone        string
	Local        string
	CurrentValue string
}

func AlarmMessages() []Message {
	var messages []Message

	smartLightData, err := fetchSmartLight()
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	userData, userError := fetchUsers()
	if userError != nil {
		log.Fatalf("Error fetching data: %v", userError)
	}
	for _, item := range userData {
		for _, alarm := range item.Alarms {
			if !alarm.AlreadyPlayed {
				messages = append(messages, Message{DEVEUI: alarm.DEVEUI, Type: alarm.Type, Trigger: alarm.Trigger, TriggerType: alarm.TriggerType, TriggerAt: alarm.TriggerAt, Phone: item.Phone, Local: alarm.Local})
			}
		}

		// fmt.Println(item.Phone)
	}

	for _, smartLight := range smartLightData {
		val := reflect.ValueOf(smartLight.Fields)
		for i := 0; i < len(messages); i++ {
			message := &messages[i]
			fieldName := message.TriggerType
			if len(fieldName) > 0 {
				fieldName = string(unicode.ToUpper(rune(fieldName[0]))) + fieldName[1:]
			}

			messageTrigger, err := strconv.ParseFloat(message.Trigger, 64)
			if err != nil {
				fmt.Printf("Error: %e", err)
			}

			fieldVal := val.FieldByName(fieldName)
			if fieldVal.IsValid() && smartLight.Tags.DeviceId == message.DEVEUI {
				triggerValue := float32(fieldVal.Float())

				if message.TriggerAt == "higher" {
					if float32(messageTrigger) > triggerValue {
						messages = append(messages[:i], messages[i+1:]...)
						i--
					}
				} else {
					if float32(messageTrigger) < triggerValue {
						messages = append(messages[:i], messages[i+1:]...)
						i--
					}
				}
				formattedString := fmt.Sprintf("%.2f", triggerValue)
				message.CurrentValue = formattedString
			}
		}
	}

	return messages
}

func main() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			// err := godotenv.Load()
			// if err != nil {
			// 	log.Fatalf("Error loading .env file: %v", err)
			// }

			accessToken := os.Getenv("ACCESS_TOKEN_META")
			phoneNumberID := os.Getenv("PHONE_NUMBER_ID_META")

			var messages []Message = AlarmMessages()

			for _, message := range messages {
				phoneNumber := "55" + message.Phone

				// POST payload
				payload := map[string]interface{}{
					"messaging_product": "whatsapp",
					"to":                phoneNumber,
					"type":              "template",
					"template": map[string]interface{}{
						"name": "alarme_smartcampus",
						"language": map[string]string{
							"code": "pt_BR",
						},
						"components": []map[string]interface{}{
							{
								"type": "body",
								"parameters": []map[string]string{
									{"type": "text", "text": message.Type},
									{"type": "text", "text": message.DEVEUI},
									{"type": "text", "text": message.TriggerType},
									{"type": "text", "text": message.CurrentValue},
									{"type": "text", "text": message.Trigger},
								},
							},
						},
					},
				}

				jsonPayload, err := json.Marshal(payload)
				if err != nil {
					log.Fatalf("Error creating JSON payload: %v", err)
				}

				// POST config
				url := fmt.Sprintf("https://graph.facebook.com/v17.0/%s/messages", phoneNumberID)
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
				if err != nil {
					log.Fatalf("Error creating request: %v", err)
				}

				// Headers
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.Header.Set("Content-Type", "application/json")

				// Request
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("Error sending request: %v", err)
				}
				// defer resp.Body.Close()

				// response
				if resp.StatusCode != http.StatusOK {
					log.Printf("Error sending message to %s. HTTP: %d\n", phoneNumber, resp.StatusCode)
				} else {
					fmt.Printf("Message sent to %s with success!\n", phoneNumber)
				}
			}

			updateAlarmAlreadyPlayedOnSupabase(messages)
		}
	}
}
