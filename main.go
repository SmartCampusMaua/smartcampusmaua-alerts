package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"crypto/tls"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
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
	Status       string `json:"status"`
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
	FirmwareVersion        int64   `json:"firmwareVersion"`
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

type SprinklerFields struct {
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
	Solenoid1      bool    `json:"solenoid1"`
	Solenoid2      bool    `json:"solenoid2"`
	Solenoid3      bool    `json:"solenoid3"`
	TxBandWidth    float64 `json:"txBandWidth"`
	TxFrequency    float64 `json:"txFrequency"`
	TxSpreadFactor float64 `json:"txSpreadFactor"`
}

type SprinklerData struct {
	Fields    SprinklerFields `json:"fields"`
	Name      string          `json:"name"`
	Tags      Tags            `json:"tags"`
	Timestamp float64         `json:"timestamp"`
}

func fetchSprinkler() ([]SprinklerData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/Sprinkler/all?interval=30"

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
	var data []SprinklerData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]SprinklerData, 0)
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

type SoilMoisture3DepthLevelsFields struct {
	BoardVoltage            float64 `json:"boardVoltage"`
	Data                    string  `json:"data"`
	FCnt                    float64 `json:"fCnt"`
	FPort                   float64 `json:"fPort"`
	RxAlt_0                 float64 `json:"rxAlt_0"`
	RxLat_0                 float64 `json:"rxLat_0"`
	RxLon_0                 float64 `json:"rxLon_0"`
	RxRssi_0                float64 `json:"rxRssi_0"`
	RxSnr_0                 float64 `json:"rxSnr_0"`
	SoilMoistureDepthLevel1 float64 `json:"soilMoistureDepthLevel1"`
	SoilMoistureDepthLevel2 float64 `json:"soilMoistureDepthLevel2"`
	SoilMoistureDepthLevel3 float64 `json:"soilMoistureDepthLevel3"`
	TxBandWidth             float64 `json:"txBandWidth"`
	TxFrequency             float64 `json:"txFrequency"`
	TxSpreadFactor          float64 `json:"txSpreadFactor"`
}

type SoilMoisture3DepthLevelsData struct {
	Fields    SoilMoisture3DepthLevelsFields `json:"fields"`
	Name      string                         `json:"name"`
	Tags      Tags                           `json:"tags"`
	Timestamp float64                        `json:"timestamp"`
}

func fetchSoilMoisture3DepthLevels() ([]SoilMoisture3DepthLevelsData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/LNS/SoilMoisture3DepthLevels/all?interval=30"

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
	var data []SoilMoisture3DepthLevelsData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]SoilMoisture3DepthLevelsData, 0)
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

type EvseStatusNotificationFields struct {
	// Info						float64 `json:"info"`
	ErrorCode string `json:"errorCode"`
}

type EvseStatusNotificationData struct {
	Fields    EvseStatusNotificationFields `json:"fields"`
	Name      string                       `json:"name"`
	Tags      Tags                         `json:"tags"`
	Timestamp float64                      `json:"timestamp"`
}

func fetchEvseStatusNotification() ([]EvseStatusNotificationData, error) {
	// API URL
	apiUrl := "https://smartcampus-k8s.maua.br/api/timeseries/v0.3/IMT/EVSE/StatusNotification/all?interval=57600" // 40 days

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

	// Parse the JSON response into a slice of EvseStatusNotificationData (since the response is an array)
	var data []EvseStatusNotificationData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Remove duplicates based on DeviceId
	uniqueData := make([]EvseStatusNotificationData, 0)
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
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
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

	alarmQuery := `SELECT "id", "userId", "type", "local", "deveui", "trigger", "triggerAt", "triggerType", "alreadyPlayed", "lastPlayed", "actionSensor" FROM "Alarms"`

	alarmRows, err := db.Query(alarmQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("error querying alarms")
	}
	defer alarmRows.Close()

	var alarmData []Alarm
	for alarmRows.Next() {
		var loc Alarm

		err := alarmRows.Scan(&loc.Id, &loc.UserId, &loc.Type, &loc.Local, &loc.DeviceId, &loc.Trigger, &loc.TriggerAt, &loc.TriggerType, &loc.AlreadyPlayed, &loc.LastPlayed, &loc.ActionSensor)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning alarm row: %v", err)
		}

		alarmData = append(alarmData, loc)
	}

	return userData, alarmData, nil
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

	query := `UPDATE "Alarms" SET "alreadyPlayed" = true WHERE "id" = $1`

	for _, message := range messages {
		_, updateAlarmErr := db.Exec(query, message.MessageAlarm.Id)
		if updateAlarmErr != nil {
			fmt.Printf("update alarm alreadyPlayed error: %v", updateAlarmErr)
		}
	}
}

type UserData struct {
	Id     int64  `json:"id"`
	Phone  string `json:"phone"`
	Alarms Alarm
}

type Alarm struct {
	Id            int64  `json:"id"`
	UserId        int64  `json:"userId"`
	Type          string `json:"type"`
	Local         string `json:"local"`
	DeviceId      string `json:"deveui"`
	Trigger       string `json:"trigger"`
	TriggerAt     string `json:"triggerAt"`
	TriggerType   string `json:"triggerType"`
	AlreadyPlayed bool   `json:"alreadyPlayed"`
	LastPlayed    string `json:"lastPlayed"`
	ActionSensor  string `json:"actionSensor"`
}

type Message struct {
	MessageAlarm Alarm
	Phone        string
	CurrentValue string
}

func AlarmMessages() []Message {
	connStr := os.Getenv("supabaseConnection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting to supabase:%v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error connecting to supabase:%v\n", err)
	}

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
	sprinklerData, err := fetchSprinkler()
	if err != nil {
		log.Fatalf("Error fetching sprinkler data: %v", err)
	}
	soilMoisture3DepthLevelsData, err := fetchSoilMoisture3DepthLevels()
	if err != nil {
		log.Fatalf("Error fetching soil moisture data: %v", err)
	}

	evseStatusNotificationData, err := fetchEvseStatusNotification()
	if err != nil {
		log.Fatalf("Error fetching evseStatusNotification data: %v", err)
	}

	userData, Alarms, userError := fetchUsers()
	if userError != nil {
		log.Fatalf("Error fetching data: %v", userError)
	}

	layout := "2006-01-02 15:04:05.999999999-07:00"
	timeNow := time.Now().Format("2006-01-02 15:04:05.999999999-07:00")
	var lastPlayedTime time.Time
	alarmTimeEnv := os.Getenv("ALARM_TIMER")
	alarmTime, _ := strconv.Atoi(alarmTimeEnv)
	for _, item := range userData {
		for _, alarm := range Alarms {
			if alarm.LastPlayed != "0" {
				lastPlayedTime, err = time.Parse(layout, alarm.LastPlayed)
			}
			if err != nil {
				fmt.Printf("Error getting lastPlayed: %v \n", err)
			}
			if alarm.AlreadyPlayed && time.Since(lastPlayedTime) >= time.Duration(alarmTime)*time.Second {
				alarm.AlreadyPlayed = false
				db.Exec(`UPDATE "Alarms" SET "alreadyPlayed" = $1 WHERE "id" = $2`, alarm.AlreadyPlayed, alarm.Id)
			}
			if !alarm.AlreadyPlayed && item.Phone != "" && item.Id == int64(alarm.UserId) {
				db.Exec(`UPDATE "Alarms" SET "lastPlayed" = $1 WHERE "id" = $2`, timeNow, alarm.Id)
				messages = append(messages, Message{MessageAlarm: alarm, Phone: item.Phone})
			}
		}
	}

	for _, message := range messages {
		dataType := message.MessageAlarm.Type
		dataTriggerType := message.MessageAlarm.TriggerType
		deviceId := message.MessageAlarm.DeviceId
		trigger, _ := strconv.ParseFloat(message.MessageAlarm.Trigger, 64)
		triggerBool, _ := strconv.ParseBool(message.MessageAlarm.TriggerAt)
		triggerAt := message.MessageAlarm.TriggerAt
		var currentValue *float64
		var currentBool *bool
		var currentString *string
		var canAddToMessages = false
		messageToSave := message

		switch dataType {
		case "SmartLight":
			{
				var dataToPass SmartLightData
				for _, smartLight := range smartLightData {
					if smartLight.Tags.DeviceId == deviceId {
						dataToPass = smartLight
					}
				}

				switch dataTriggerType {
				case "batteryVoltage":
					{
						currentValue = &dataToPass.Fields.BatteryVoltage
					}
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "humidity":
					{
						currentValue = &dataToPass.Fields.Humidity
					}
				case "luminosity":
					{
						currentValue = &dataToPass.Fields.Luminosity
					}
				case "movement":
					{
						currentValue = &dataToPass.Fields.Movement
					}
				case "temperature":
					{
						currentValue = &dataToPass.Fields.Temperature
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				}
			}

		case "WaterTankLevel":
			{
				var dataToPass WaterTankData
				for _, waterTank := range waterTankData {
					if waterTank.Tags.DeviceId == deviceId {
						dataToPass = waterTank
					}
				}

				switch dataTriggerType {
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "distance":
					{
						currentValue = &dataToPass.Fields.Distance
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				}
			}

		case "Hydrometer":
			{
				var dataToPass HydrometerData
				for _, hydrometer := range hydrometerData {
					if hydrometer.Tags.DeviceId == deviceId {
						dataToPass = hydrometer
					}
				}

				switch dataTriggerType {
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "counter":
					{
						currentValue = &dataToPass.Fields.Counter
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				}
			}

		case "EnergyMeter":
			{
				var dataToPass EnergyMeterData
				for _, energyMeter := range energyMeterData {
					if energyMeter.Tags.DeviceId == deviceId {
						dataToPass = energyMeter
					}
				}

				switch dataTriggerType {
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "forwardEnergy":
					{
						currentValue = &dataToPass.Fields.ForwardEnergy
					}
				case "reverseEnergy":
					{
						currentValue = &dataToPass.Fields.ReverseEnergy
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				}
			}

		case "WeatherStation":
			{
				var dataToPass WeatherStationData
				for _, weatherStation := range weatherStationData {
					if weatherStation.Tags.DeviceId == deviceId {
						dataToPass = weatherStation
					}
				}

				switch dataTriggerType {
				case "c1Count":
					{
						currentValue = &dataToPass.Fields.C1Count
					}
				case "c1State":
					{
						currentBool = &dataToPass.Fields.C1State
					}
				case "c2Count":
					{
						currentValue = &dataToPass.Fields.C2Count
					}
				case "c2State":
					{
						currentBool = &dataToPass.Fields.C2State
					}
				case "emwAtmPres":
					{
						currentValue = &dataToPass.Fields.EmwAtmPres
					}
				case "emwAvgWindSpeed":
					{
						currentValue = &dataToPass.Fields.EmwAvgWindSpeed
					}
				case "emwGustWindSpeed":
					{
						currentValue = &dataToPass.Fields.EmwGustWindSpeed
					}
				case "emwHumidity":
					{
						currentValue = &dataToPass.Fields.EmwHumidity
					}
				case "emwLuminosity":
					{
						currentValue = &dataToPass.Fields.EmwLuminosity
					}
				case "emwRainLevel":
					{
						currentValue = &dataToPass.Fields.EmwRainLevel
					}
				case "emwSolarRadiation":
					{
						currentValue = &dataToPass.Fields.EmwSolarRadiation
					}
				case "emwTemperature":
					{
						currentValue = &dataToPass.Fields.EmwTemperature
					}
				case "emwUv":
					{
						currentValue = &dataToPass.Fields.EmwUv
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				} else if triggerAt == "true" && triggerBool && *currentBool {
					canAddToMessages = true
				} else if triggerAt == "false" && !triggerBool && !*currentBool {
					canAddToMessages = true
				}
			}

		case "Sprinkler":
			{
				var dataToPass SprinklerData
				for _, sprinkler := range sprinklerData {
					if sprinkler.Tags.DeviceId == deviceId {
						dataToPass = sprinkler
					}
				}

				switch dataTriggerType {
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "counter":
					{
						currentValue = &dataToPass.Fields.Counter
					}
				case "solenoid1":
					{
						currentBool = &dataToPass.Fields.Solenoid1
					}
				case "solenoid2":
					{
						currentBool = &dataToPass.Fields.Solenoid2
					}
				case "solenoid3":
					{
						currentBool = &dataToPass.Fields.Solenoid3
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				} else if triggerAt == "true" && triggerBool && *currentBool {
					canAddToMessages = true
				} else if triggerAt == "false" && !triggerBool && !*currentBool {
					canAddToMessages = true
				}
			}

		case "SoilMoisture3DepthLevels":
			{
				var dataToPass SoilMoisture3DepthLevelsData
				for _, soilMoisture := range soilMoisture3DepthLevelsData {
					if soilMoisture.Tags.DeviceId == deviceId {
						soilMoisture.Fields.SoilMoistureDepthLevel1 /= 100
						soilMoisture.Fields.SoilMoistureDepthLevel2 /= 100
						soilMoisture.Fields.SoilMoistureDepthLevel3 /= 100
						dataToPass = soilMoisture
					}
				}

				switch dataTriggerType {
				case "boardVoltage":
					{
						currentValue = &dataToPass.Fields.BoardVoltage
					}
				case "soilMoistureDepthLevel1":
					{
						currentValue = &dataToPass.Fields.SoilMoistureDepthLevel1
					}
				case "soilMoistureDepthLevel2":
					{
						currentValue = &dataToPass.Fields.SoilMoistureDepthLevel2
					}
				case "soilMoistureDepthLevel3":
					{
						currentValue = &dataToPass.Fields.SoilMoistureDepthLevel3
					}
				}
				if triggerAt == "higher" && trigger < *currentValue {
					canAddToMessages = true
				} else if triggerAt == "lower" && trigger > *currentValue {
					canAddToMessages = true
				}
			}

		// 		case "stopTime":
		// 			{
		// 				currentValue = &dataToPass.Fields.StopTime
		// 			}
		// 		}
		// 	}
		case "Evse":
			{
				var dataToPass EvseStatusNotificationData
				for _, evseStatusNotification := range evseStatusNotificationData {
					if evseStatusNotification.Tags.DeviceId == deviceId {
						dataToPass = evseStatusNotification
					}
				}

				switch dataTriggerType {
				case "status":
					{
						currentString = &dataToPass.Tags.Status
					}
				}
				if *currentString == "Available" {
					canAddToMessages = true
				}
			}
		}

		if canAddToMessages {
			if currentValue != nil {
				messageToSave.CurrentValue = fmt.Sprintf("%v", *currentValue)
			} else if currentBool != nil {
				if *currentBool {
					messageToSave.CurrentValue = "Verdadeiro"
				} else {
					messageToSave.CurrentValue = "Falso"
				}
			} else {
				messageToSave.CurrentValue = "Disponível"
			}
			finalMessages = append(finalMessages, messageToSave)
		}
	}
	// Places new lastPlayed value on database
	for _, finalMessage := range finalMessages {
		if finalMessage.MessageAlarm.LastPlayed == "0" {
			_, updateAlarmErr := db.Exec(`UPDATE "Alarms" SET "lastPlayed" = $1 WHERE "id" = $2`, timeNow, finalMessage.MessageAlarm.Id)
			if updateAlarmErr != nil {
				fmt.Printf("update alarm alreadyPlayed error: %v", updateAlarmErr)
			}
		}
		// Alarm_History always has alreadyPlayed as false due to the order of inserts, may need to fix in the future if history needs change
		_, alarmsHistoryErr := db.Exec(`INSERT INTO "Alarms_History" ("id", "userId", "type", "local", "deveui", "trigger", "triggerAt", "triggerType", "lastPlayed", "currentValue", "actionSensor") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, finalMessage.MessageAlarm.Id, finalMessage.MessageAlarm.UserId, finalMessage.MessageAlarm.Type, finalMessage.MessageAlarm.Local, finalMessage.MessageAlarm.DeviceId, finalMessage.MessageAlarm.Trigger, finalMessage.MessageAlarm.TriggerAt, finalMessage.MessageAlarm.TriggerType, timeNow, finalMessage.CurrentValue, finalMessage.MessageAlarm.ActionSensor)
		if alarmsHistoryErr != nil {
			fmt.Printf("insert alarmHistory error: %v\n", alarmsHistoryErr)
		}
	}

	return finalMessages
}

func main() {
	tickerTimeEnv := os.Getenv("TICKER_TIME")
	tickerTime, _ := strconv.Atoi(tickerTimeEnv)
	ticker := time.NewTicker(time.Duration(tickerTime) * time.Second)
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
				var payload map[string]interface{}

				if message.MessageAlarm.Type == "Evse" {
					// POST payload
					payload = map[string]interface{}{
						"messaging_product": "whatsapp",
						"to":                phoneNumber,
						"type":              "template",
						"template": map[string]interface{}{
							"name": "evse",
							"language": map[string]string{
								"code": "pt_BR",
							},
							"components": []map[string]interface{}{
								{
									"type": "body",
									"parameters": []map[string]string{
										{"type": "text", "text": message.MessageAlarm.Type},
										{"type": "text", "text": message.MessageAlarm.DeviceId},
										{"type": "text", "text": message.MessageAlarm.Local},
									},
								},
							},
						},
					}
				} else if message.CurrentValue == "Verdadeiro" || message.CurrentValue == "Falso" {
					// POST payload
					payload = map[string]interface{}{
						"messaging_product": "whatsapp",
						"to":                phoneNumber,
						"type":              "template",
						"template": map[string]interface{}{
							"name": "alarme_bool",
							"language": map[string]string{
								"code": "pt_BR",
							},
							"components": []map[string]interface{}{
								{
									"type": "body",
									"parameters": []map[string]string{
										{"type": "text", "text": message.MessageAlarm.Type},
										{"type": "text", "text": message.MessageAlarm.DeviceId},
										{"type": "text", "text": message.MessageAlarm.Local},
										{"type": "text", "text": message.MessageAlarm.TriggerType},
										{"type": "text", "text": message.CurrentValue},
										{"type": "text", "text": message.CurrentValue},
									},
								},
							},
						},
					}
				} else {
					var triggerAt string
					if message.MessageAlarm.TriggerAt == "higher" {
						triggerAt = "acima"
					} else {
						triggerAt = "abaixo"
					}

					// POST payload
					payload = map[string]interface{}{
						"messaging_product": "whatsapp",
						"to":                phoneNumber,
						"type":              "template",
						"template": map[string]interface{}{
							"name": "alerta_smart",
							"language": map[string]string{
								"code": "pt_BR",
							},
							"components": []map[string]interface{}{
								{
									"type": "body",
									"parameters": []map[string]string{
										{"type": "text", "text": message.MessageAlarm.Type},
										{"type": "text", "text": message.MessageAlarm.DeviceId},
										{"type": "text", "text": message.MessageAlarm.Local},
										{"type": "text", "text": message.MessageAlarm.TriggerType},
										{"type": "text", "text": triggerAt},
										{"type": "text", "text": message.CurrentValue},
										{"type": "text", "text": message.MessageAlarm.Trigger},
									},
								},
							},
						},
					}
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

				// Actuators
				var action string
				if message.MessageAlarm.ActionSensor == "sprinklersOn" {
					action = "AFQ="
				} else if message.MessageAlarm.ActionSensor == "sprinklersOff" {
					action = "AKg="
				} else {
					action = ""
				}

				payload = map[string]interface{}{
					"application": "SmartLight",          // Application Name registered in the corresponding NetworkServer
					"etc":         "imt",                 // NetworkServer to be queued
					"reference":   "test-node-red",       // Reference
					"deviceId":    "0004a30b00e94314",    // Device ID
					"confirmed":   false,                 // Confirmed flag
					"fPort":       100,                   // lora downlink fPort
					"data":        action,                // Downlink data
					"timestamp":   time.Now().UnixNano(), // Current timestamp in nanoseconds
				}

				jsonData, err := json.Marshal(payload)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}

				url = "https://smartcampus-k8s.maua.br/api/ingestion/v0.1/IMT/LNS/Command/all"

				req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return
				}

				req.Header.Set("Content-Type", "application/json")

				client = &http.Client{}
				resp, err = client.Do(req)
				if err != nil {
					fmt.Println("Error making request:", err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					var responseData map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
						fmt.Println("Error decoding response:", err)
						return
					}
					fmt.Println("Response:", responseData)
					fmt.Println("Post OK!")
				} else {
					fmt.Printf("HTTP Error: %d\n", resp.StatusCode)
				}

				// Alert
				var deviceType string
				var measurement string
				if message.MessageAlarm.Type == "Evse" {
					deviceType = "EVSE"
					measurement = "MeterValues"
				} else {
					deviceType = "LNS"
					measurement = message.MessageAlarm.Type
				}

				payload = map[string]interface{}{
					"deviceType":   deviceType,                    // LNS, EVSE
					"measurement":  measurement,                   // SmartLight, WeatherStation
					"deviceId":     message.MessageAlarm.DeviceId, // Device ID
					"trigger":      message.MessageAlarm.Trigger,
					"triggerAt":    message.MessageAlarm.TriggerAt,
					"triggerType":  message.MessageAlarm.TriggerType,
					"lastPlayed":   message.MessageAlarm.LastPlayed,
					"actionSensor": message.MessageAlarm.ActionSensor,
					"currentValue": message.CurrentValue,
					"data":         "Alert SmartCampus",   // Downlink data
					"timestamp":    time.Now().UnixNano(), // Current timestamp in nanoseconds
					"etc":          "imt",                 // NetworkServer to be queued
				}

				jsonData, err = json.Marshal(payload)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}

				var sbUrl strings.Builder
				sbUrl.WriteString(`https://smartcampus-k8s.maua.br/api/ingestion/v0.1/IMT/`)
				sbUrl.WriteString(deviceType)
				sbUrl.WriteString(`/Alert/all`)

				url = sbUrl.String()

				req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return
				}

				req.Header.Set("Content-Type", "application/json")

				client = &http.Client{}
				resp, err = client.Do(req)
				if err != nil {
					fmt.Println("Error making request:", err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					var responseData map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
						fmt.Println("Error decoding response:", err)
						return
					}
					fmt.Println("Response:", responseData)
					fmt.Println("Post OK!")
				} else {
					fmt.Printf("HTTP Error: %d\n", resp.StatusCode)
				}
			}
			updateAlarmAlreadyPlayedOnSupabase(messages)
		}
	}
}
