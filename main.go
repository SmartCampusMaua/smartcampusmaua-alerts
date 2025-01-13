package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unicode"

	"crypto/tls"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type Tags struct {
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

type SmartLightData struct {
	Fields    SmartLightFields `json:"fields"`
	Name      string           `json:"name"`
	Tags      Tags             `json:"tags"`
	Timestamp float64          `json:"timestamp"`
}

func fetchSmartLight() ([]SmartLightData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/SmartLight/all?interval=30"

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

type WaterTankFields struct {
	BoardVoltage   float64 `json:"boardVoltage"`
	Data           string  `json:"data"`
	Distance       float64 `json:"distance"`
	FCnt           float64 `json:"fCnt"`
	FPort          float64 `json:"fPort"`
	RxAlt_0        float64 `json:"rxAlt_0"`
	RxLat_0        float64 `json:"rxLat_0"`
	RxLon_0        float64 `json:"rxLon_0"`
	RxRssi_0       float64 `json:"rxRssi_0"`
	RxSnr_0        float64 `json:"rxSnr_0"`
	TxBandWidth    float64 `json:"txBandWidth"`
	TxFrequency    float64 `json:"txFrequency"`
	TxSpreadFactor float64 `json:"txSpreadFactor"`
}

type WaterTankData struct {
	Fields    WaterTankFields `json:"fields"`
	Name      string          `json:"name"`
	Tags      Tags            `json:"tags"`
	Timestamp float64         `json:"timestamp"`
}

func fetchWaterTank() ([]WaterTankData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/WaterTankLevel/all?interval=30"

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
	var data []WaterTankData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]WaterTankData, 0)
	seenDevices := make(map[string]bool)

	for _, item := range data {
		if !seenDevices[item.Tags.DeviceId] {
			uniqueData = append(uniqueData, item)
			seenDevices[item.Tags.DeviceId] = true
		}
	}

	// Return the filtered data
	return uniqueData, nil
}

type HydrometerFields struct {
	BoardVoltage   float64 `json:"boardVoltage"`
	Counter        float64 `json:"counter"`
	Data           string  `json:"data"`
	FCnt           float64 `json:"fCnt"`
	FPort          float64 `json:"fPort"`
	RxAlt_0        float64 `json:"rxAlt_0"`
	RxLat_0        float64 `json:"rxLat_0"`
	RxLon_0        float64 `json:"rxLon_0"`
	RxRssi_0       float64 `json:"rxRssi_0"`
	RxSnr_0        float64 `json:"rxSnr_0"`
	TxBandWidth    float64 `json:"txBandWidth"`
	TxFrequency    float64 `json:"txFrequency"`
	TxSpreadFactor float64 `json:"txSpreadFactor"`
}

type HydrometerData struct {
	Fields    HydrometerFields `json:"fields"`
	Name      string           `json:"name"`
	Tags      Tags             `json:"tags"`
	Timestamp float64          `json:"timestamp"`
}

func fetchHydrometer() ([]HydrometerData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/Hydrometer/all?interval=30"

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
	var data []HydrometerData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]HydrometerData, 0)
	seenDevices := make(map[string]bool)

	for _, item := range data {
		if !seenDevices[item.Tags.DeviceId] {
			uniqueData = append(uniqueData, item)
			seenDevices[item.Tags.DeviceId] = true
		}
	}

	// Return the filtered data
	return uniqueData, nil
}

type EnergyMeterFields struct {
	BoardVoltage   float64 `json:"boardVoltage"`
	Data           string  `json:"data"`
	FCnt           float64 `json:"fCnt"`
	FPort          float64 `json:"fPort"`
	ForwardEnergy  float64 `json:"forwardEnergy"`
	ReverseEnergy  float64 `json:"reverseEnergy"`
	RxAlt_0        float64 `json:"rxAlt_0"`
	RxLat_0        float64 `json:"rxLat_0"`
	RxLon_0        float64 `json:"rxLon_0"`
	RxRssi_0       float64 `json:"rxRssi_0"`
	RxSnr_0        float64 `json:"rxSnr_0"`
	TxBandWidth    float64 `json:"txBandWidth"`
	TxFrequency    float64 `json:"txFrequency"`
	TxSpreadFactor float64 `json:"txSpreadFactor"`
}

type EnergyMeterData struct {
	Fields    EnergyMeterFields `json:"fields"`
	Name      string            `json:"name"`
	Tags      Tags              `json:"tags"`
	Timestamp float64           `json:"timestamp"`
}

func fetchEnergyMeter() ([]EnergyMeterData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/EnergyMeter/all?interval=30"

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
	var data []EnergyMeterData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]EnergyMeterData, 0)
	seenDevices := make(map[string]bool)

	for _, item := range data {
		if !seenDevices[item.Tags.DeviceId] {
			uniqueData = append(uniqueData, item)
			seenDevices[item.Tags.DeviceId] = true
		}
	}

	// Return the filtered data
	return uniqueData, nil
}

type WeatherStationFields struct {
	C1Count                float64 `json:"c1Count"`
	C1State                bool    `json:"c1State"`
	C2Count                float64 `json:"c2Count"`
	C2State                bool    `json:"c2State"`
	Data                   string  `json:"data"`
	EmwAtmPres             float64 `json:"emwAtmPres"`
	EmwAvgWindSpeed        float64 `json:"emwAvgWindSpeed"`
	EmwGustWindSpeed       float64 `json:"emwGustWindSpeed"`
	EmwHumidity            float64 `json:"emwHumidity"`
	EmwLuminosity          float64 `json:"emwLuminosity"`
	EmwRainLevel           float64 `json:"emwRainLevel"`
	EmwSolarRadiation      float64 `json:"emwSolarRadiation"`
	EmwTemperature         float64 `json:"emwTemperature"`
	EmwUv                  float64 `json:"emwUv"`
	EmwWindDirection       float64 `json:"emwWindDirection"`
	EnvSensorFailStatus    bool    `json:"envSensorFailStatus"`
	FCnt                   float64 `json:"fCnt"`
	FPort                  float64 `json:"fPort"`
	FirmwareVersion        uint64
	InternalBatteryVoltage float64 `json:"internalBatteryVoltage"`
	InternalHumidity       float64 `json:"internalHumidity"`
	InternalTemperature    float64 `json:"internalTemperature"`
	PowerSource            bool    `json:"powerSource"`
	RxAlt_0                float64 `json:"rxAlt_0"`
	RxLat_0                float64 `json:"rxLat_0"`
	RxLon_0                float64 `json:"rxLon_0"`
	RxRssi_0               float64 `json:"rxRssi_0"`
	RxSnr_0                float64 `json:"rxSnr_0"`
	TxBandWidth            float64 `json:"txBandWidth"`
	TxFrequency            float64 `json:"txFrequency"`
	TxSpreadFactor         float64 `json:"txSpreadFactor"`
}

type WeatherStationData struct {
	Fields    WeatherStationFields `json:"fields"`
	Name      string               `json:"name"`
	Tags      Tags                 `json:"tags"`
	Timestamp float64              `json:"timestamp"`
}

func fetchWeatherStation() ([]WeatherStationData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/WeatherStation/all?interval=30"

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
	var data []WeatherStationData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]WeatherStationData, 0)
	seenDevices := make(map[string]bool)

	for _, item := range data {
		if !seenDevices[item.Tags.DeviceId] {
			uniqueData = append(uniqueData, item)
			seenDevices[item.Tags.DeviceId] = true
		}
	}

	// Return the filtered data
	return uniqueData, nil
}

func fetchUsers() ([]UserData, []Alarm, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("supabaseConnection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	userQuery := `SELECT "id", "phone" FROM "User"`

	userRows, err := db.Query(userQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("error querying users: %v", err)
	}
	defer userRows.Close()

	var userData []UserData
	for userRows.Next() {
		var loc UserData

		err := userRows.Scan(&loc.Id, &loc.Phone)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning user row: %v", err)
		}

		userData = append(userData, loc)
	}

	alarmQuery := fmt.Sprintf(`SELECT "id", "type", "local", "deveui", "trigger", "triggerAt", "triggerType", "alreadyPlayed" FROM "Alarms"`)

	alarmRows, err := db.Query(alarmQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("error querying alarms")
	}
	defer alarmRows.Close()

	var alarmData []Alarm
	for alarmRows.Next() {
		var loc Alarm

		err := alarmRows.Scan(&loc.Id, &loc.Type, &loc.Local, &loc.Deveui, &loc.Trigger, &loc.TriggerAt, &loc.TriggerType, &loc.AlreadyPlayed)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning alarm row: %v", err)
		}

		alarmData = append(alarmData, loc)
	}

	return userData, alarmData, nil
}

func updateAlarmAlreadyPlayedOnSupabase(messages []Message) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	Id     int64  `json:"id"`
	Phone  string `json:"phone"`
	Alarms Alarm
}

type Alarm struct {
	Id            int8   `json:"id"`
	UserId        int8   `json:"userId"`
	Type          string `json:"type"`
	Local         string `json:"local"`
	Deveui        string `json:"deveui"`
	Trigger       string `json:"trigger"`
	TriggerAt     string `json:"triggerAt"`
	TriggerType   string `json:"triggerType"`
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

func getFieldValue(data interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s does not exist", fieldName)
	}

	switch field.Kind() {
	case reflect.Float64:
		return field.Float(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int(), nil
	case reflect.String:
		return field.String(), nil
	case reflect.Bool:
		return field.Bool(), nil
	default:
		return nil, fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

func AlarmMessages() []Message {
	var messages []Message
	var finalMessages []Message

	smartLightData, err := fetchSmartLight()
	if err != nil {
		log.Fatalf("Error fetching smartlight data: %v", err)
	}
	waterTankData, err := fetchWaterTank()
	if err != nil {
		log.Fatalf("Error fetching watertank data: %v", err)
	}
	hydrometerData, err := fetchHydrometer()
	if err != nil {
		log.Fatalf("Error fetching hydrometer data: %v", err)
	}
	energyMeterData, err := fetchEnergyMeter()
	if err != nil {
		log.Fatalf("Error fetching energymeter data: %v", err)
	}
	weatherStationData, err := fetchWeatherStation()
	if err != nil {
		log.Fatalf("Error fetching weatherstation data: %v", err)
	}

	userData, Alarms, userError := fetchUsers()
	if userError != nil {
		log.Fatalf("Error fetching data: %v", userError)
	}
	for _, item := range userData {
		for _, alarm := range Alarms {
			if !alarm.AlreadyPlayed {
				messages = append(messages, Message{DEVEUI: alarm.Deveui, Type: alarm.Type, Trigger: alarm.Trigger, TriggerType: alarm.TriggerType, TriggerAt: alarm.TriggerAt, Phone: item.Phone, Local: alarm.Local})
			}
		}
	}

	for _, message := range messages {
		dataType := message.Type
		dataTriggerType := message.TriggerType
		var dataValue interface{}
		var err error

		switch dataTriggerType {
		case "SmartLight":
			dataValue, err = getFieldValue(smartLightData, dataType)
		case "WaterTankLevel":
			dataValue, err = getFieldValue(waterTankData, dataType)
		case "Hydrometer":
			dataValue, err = getFieldValue(hydrometerData, dataType)
		case "EnergyMeter":
			dataValue, err = getFieldValue(energyMeterData, dataType)
		case "WeatherStation":
			dataValue, err = getFieldValue(weatherStationData, dataType)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		switch v := dataValue.(type) {
		case float64:
      triggerValue, _ := strconv.ParseFloat(message.Trigger, 64)
			if message.TriggerAt == "higher" && triggerValue > v {
				finalMessages = append(finalMessages, message)
			} else if message.TriggerAt == "lower" && triggerValue < v {
				finalMessages = append(finalMessages, message)
			}

			message.CurrentValue = fmt.Sprintf("%v", v)
		case int64:
      triggerValue, _ := strconv.ParseInt(message.Trigger, 10, 64)
			if message.TriggerAt == "higher" && triggerValue > v {
				finalMessages = append(finalMessages, message)
			} else if message.TriggerAt == "lower" && triggerValue < v {
				finalMessages = append(finalMessages, message)
			}

			message.CurrentValue = fmt.Sprintf("%v", v)
		case string:
			fmt.Printf("String value: %s\n", v) // Tratar depois, não sei como será o caso string então não adianta fazer agora
		case bool:
			if message.TriggerAt == "true" && v == true {
				finalMessages = append(finalMessages, message)
			} else if message.Trigger == "false" && v == false {
				finalMessages = append(finalMessages, message)
			}

			message.CurrentValue = fmt.Sprintf("%v", v)
		}
	}

	return finalMessages
}

func main() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file: %v", err)
			}

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

			// updateAlarmAlreadyPlayedOnSupabase(messages)
		}
	}
}
