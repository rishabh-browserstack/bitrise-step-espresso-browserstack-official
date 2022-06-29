package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TestMapping struct {
	Name     string   `json:"name,omitempty"`
	Strategy string   `json:"strategy,omitempty"`
	Values   []string `json:"values,omitempty"`
}

type TestSharding struct {
	NumberOfShards      int           `json:"numberOfShards,omitempty"`
	Mapping             []TestMapping `json:"mapping,omitempty"`
	AutoStrategyDevices []string      `json:"devices,omitempty"`
}

type TestCases struct {
	Count  int                    `json:"count,omitempty"`
	Status map[string]interface{} `json:"status,omitempty"`
}

type Sessions struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	TestCases string `json:"start_time"`
}

// const a = {"numberOfShards": 2}, "devices": ["Google Pixel 3-9.0"]
// Input for package strategy:
// {"numberOfShards": 2, "mapping": [{"name": "Shard 1", "strategy": "package", "values": ["com.foo.login", "com.foo.logout"]}, {"name": "Shard 2", "strategy": "package", "values": ["com.foo.dashboard"]}]}
// Input for class strategy:
// {"numberOfShards": 2, "mapping": [{"name": "Shard 1", "strategy": "class", "values": ["com.foo.login.user", "com.foo.login.admin"]}, {"name": "Shard 2", "strategy": "class", "values": ["com.foo.logout.user"]}]}

type BrowserStackPayload struct {
	App                    string      `json:"app"`
	TestSuite              string      `json:"testSuite"`
	Devices                []string    `json:"devices"`
	InstrumentationLogs    bool        `json:"instrumentationLogs"`
	NetworkLogs            bool        `json:"networkLogs"`
	DeviceLogs             bool        `json:"deviceLogs"`
	DebugScreenshots       bool        `json:"debugscreenshots,omitempty"`
	VideoRecording         bool        `json:"video"`
	Project                string      `json:"project,omitempty"`
	ProjectNotifyURL       string      `json:"projectNotifyURL,omitempty"`
	UseLocal               bool        `json:"local,omitempty"`
	ClearAppData           bool        `json:"clearPackageData,omitempty"`
	SingleRunnerInvocation bool        `json:"singleRunnerInvocation,omitempty"`
	Class                  []string    `json:"class,omitempty"`
	Package                []string    `json:"package,omitempty"`
	Annotation             []string    `json:"annotation,omitempty"`
	Size                   []string    `json:"size,omitempty"`
	UseMockServer          bool        `json:"allowDeviceMockServer,omitempty"`
	UseTestSharding        interface{} `json:"shards,omitempty"`

	// non ui fields
	// EnableSpoonFramework  bool     `json:"enableSpoonFramework,omitempty"`
	// GpsLocation           string   `json:"gpsLocation,omitempty"`
	// GeoLocation           string   `json:"geoLocation,omitempty"`
	// CallbackURL           string   `json:"callbackURL,omitempty"`
	// NetworkProfile        string   `json:"networkProfile,omitempty"`
	// CustomNetwork         string   `json:"customNetwork,omitempty"`
	// Language              string   `json:"language,omitempty"`
	// Locale                string   `json:"locale,omitempty"`
	// AppStoreConfiguration string   `json:"appStoreConfiguration,omitempty"`
	// DeviceOrientation     string   `json:"deviceOrientation,omitempty"`
	// AcceptInsecureCerts   bool     `json:"acceptInsecureCerts,omitempty"`
	// UploadMedia           []string `json:"UploadMedia,omitempty"`
	// LocalIdentifier       string   `json:"localIdentifier,omitempty"`
	// IdleTimeout           string   `json:"idleTimeout,omitempty"`
}

func getDevices() ([]string, error) {
	var devices []string

	devices_input := os.Getenv("devices_list")

	if devices_input == "" {
		return devices, errors.New(fmt.Sprintf(BUILD_FAILED_ERROR, "invalid device format"))
	}

	scanner := bufio.NewScanner(strings.NewReader(devices_input))

	for scanner.Scan() {
		device := scanner.Text()
		device = strings.TrimSpace(device)

		if device == "" {
			continue
		}

		devices = append(devices, device)

	}
	return devices, nil
}

func appendExtraCapabilities(payload string) []byte {

	out := map[string]interface{}{}

	json.Unmarshal([]byte(payload), &out)

	scanner := bufio.NewScanner(strings.NewReader(os.Getenv("api_params")))
	for scanner.Scan() {
		test_sharding := scanner.Text()

		test_sharding = strings.TrimSpace(test_sharding)

		if test_sharding == "" {
			continue
		}

		test_values := strings.Split(test_sharding, "=")

		key := test_values[0]

		out[key] = test_values[1]

	}

	outputJSON, _ := json.Marshal(out)
	return outputJSON

}

func getTestFilters(payload *BrowserStackPayload) {

	scanner := bufio.NewScanner(strings.NewReader(os.Getenv("filter_test")))
	for scanner.Scan() {
		test_sharding := scanner.Text()

		test_sharding = strings.TrimSpace(test_sharding)

		if test_sharding == "" {
			continue
		}

		test_values := strings.Split(test_sharding, ",")

		for i := 0; i < len(test_values); i++ {

			test_value := strings.Split(test_values[i], " ")
			switch test_value[0] {
			case "class":
				*&payload.Class = append(*&payload.Class, test_value[1])
			case "package":
				*&payload.Package = append(*&payload.Package, test_value[1])
			case "annotation":
				*&payload.Annotation = append(*&payload.Annotation, test_value[1])
			case "size":
				*&payload.Size = append(*&payload.Size, test_value[1])
			}
		}

	}

}

func createBuildPayload() BrowserStackPayload {
	instrumentation_logs, _ := strconv.ParseBool(os.Getenv("instrumentation_logs"))
	network_logs, _ := strconv.ParseBool(os.Getenv("network_logs"))
	device_logs, _ := strconv.ParseBool(os.Getenv("device_logs"))
	debug_screenshots, _ := strconv.ParseBool(os.Getenv("debug_screenshots"))
	video_recording, _ := strconv.ParseBool(os.Getenv("video_recording"))
	use_local, _ := strconv.ParseBool(os.Getenv("use_local"))
	clear_app_data, _ := strconv.ParseBool(os.Getenv("clear_app_data"))
	use_single_runner_invocation, _ := strconv.ParseBool(os.Getenv("use_single_runner_invocation"))
	use_mock_server, _ := strconv.ParseBool(os.Getenv("use_mock_server"))

	sharding_data := TestSharding{}
	if os.Getenv("use_test_sharding") != "" {
		err := json.Unmarshal([]byte(os.Getenv("use_test_sharding")), &sharding_data)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	payload := BrowserStackPayload{
		InstrumentationLogs:    instrumentation_logs,
		NetworkLogs:            network_logs,
		DeviceLogs:             device_logs,
		DebugScreenshots:       debug_screenshots,
		VideoRecording:         video_recording,
		SingleRunnerInvocation: use_single_runner_invocation,
		Project:                os.Getenv("project"),
		ProjectNotifyURL:       os.Getenv("project_notify_url"),
		UseLocal:               use_local,
		ClearAppData:           clear_app_data,
		UseMockServer:          use_mock_server,
	}

	getTestFilters(&payload)

	if len(sharding_data.Mapping) != 0 && sharding_data.NumberOfShards != 0 {
		payload.UseTestSharding = sharding_data
	}

	payload.Devices, _ = getDevices()

	return payload
}

func failf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
	os.Exit(1)
}

func setInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				// return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear

}

func jsonParse(base64String string) map[string]interface{} {
	parsed_json := make(map[string]interface{})

	err := json.Unmarshal([]byte(base64String), &parsed_json)

	if err != nil {
		failf("Unable to parse app_upload API response: %s", err)
	}

	return parsed_json
}

func printBuildStatus(build_details map[string]interface{}) {

	log.Println("Build finished")
	log.Println("Test results summary:")

	devices := build_details["devices"].([]interface{})
	build_id := build_details["id"]
	log.Print("len", len(devices))

	if len(devices) == 1 {
		log.Print("heree")
		sessions := devices[0].(map[string]interface{})["sessions"].([]interface{})[0].(map[string]interface{})

		session_status := sessions["status"].(string)
		session_test_cases := sessions["testcases"].(map[string]interface{})
		session_test_status := session_test_cases["status"].(map[string]interface{})

		total_test := session_test_cases["count"]
		passed_test := session_test_status["passed"]
		device_name := devices[0].(map[string]interface{})["device"].(string)

		log.Print("Build Id                                            Devices                                            Status")
		log.Println("")
		if session_status == "passed" {
			log.Printf("%s                                            %s                                            PASSED (%v/%v passed)", build_id, device_name, passed_test, total_test)

		}

		if session_status == "failed" || session_status == "error" {
			log.Printf("%s                                            %s                                            FAILED (%v/%v passed)", build_id, device_name, passed_test, total_test)
		}

	} else {
		for i := 0; i < len(devices); i++ {
			sessions := devices[i].(map[string]interface{})["sessions"].([]interface{})[0].(map[string]interface{})

			session_status := sessions["status"].(string)
			session_test_cases := sessions["testcases"].(map[string]interface{})
			session_test_status := session_test_cases["status"].(map[string]interface{})

			total_test := session_test_cases["count"]
			passed_test := session_test_status["passed"]
			device_name := devices[i].(map[string]interface{})["device"].(string)

			// log.Print("Payload -> ", session_status, session_test_cases, total_test, session_status)

			log.Print("Build Id                                            Devices                                            Status")
			if session_status == "passed" {
				log.Printf("%s                %s                PASSED (%v/%v passed)", build_id, device_name, passed_test, total_test)

			}

			if session_status == "failed" || session_status == "error" {
				log.Printf("%s                %s                FAILED (%v/%v passed)", build_id, device_name, passed_test, total_test)
			}
		}
	}
}
