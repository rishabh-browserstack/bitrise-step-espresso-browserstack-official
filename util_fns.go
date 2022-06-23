package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type BrowserStackPayload struct {
	App                    string      `json:"app"`
	TestSuite              string      `json:"testSuite"`
	Devices                []string    `json:"devices"`
	InstrumentationLogs    bool        `json:"instrumentationLogs"`
	NetworkLogs            bool        `json:"networkLogs"`
	DeviceLogs             bool        `json:"deviceLogs"`
	Screenshots            bool        `json:"screenshots"`
	VideoRecording         bool        `json:"video"`
	Project                string      `json:"project,omitempty"`
	ProjectNotifyURL       string      `json:"projectNotifyURL,omitempty"`
	UseLocal               bool        `json:"useLocal,omitempty"`
	ClearAppData           bool        `json:"clearPackageData,omitempty"`
	SingleRunnerInvocation bool        `json:"singleRunnerInvocation,omitempty"`
	Class                  interface{} `json:"class,omitempty"`
	Package                []string    `json:"package,omitempty"`
	Annotation             []string    `json:"annotation,omitempty"`
	Size                   []string    `json:"size,omitempty"`
	UseMockServer          bool        `json:"allowDeviceMockServer,omitempty"`
	// UseTestSharding        interface{} `json:"shards,omitempty"`
	// Filter_tests           interface{} `json:"shards,omitempty"`
}

func createBuildPayload() BrowserStackPayload {
	instrumentation_logs, _ := strconv.ParseBool(os.Getenv("instrumentation_logs"))
	network_logs, _ := strconv.ParseBool(os.Getenv("network_logs"))
	device_logs, _ := strconv.ParseBool(os.Getenv("device_logs"))
	screenshots, _ := strconv.ParseBool(os.Getenv("screenshots"))
	video_recording, _ := strconv.ParseBool(os.Getenv("video_recording"))
	use_local, _ := strconv.ParseBool(os.Getenv("use_local"))
	clear_app_data, _ := strconv.ParseBool(os.Getenv("clear_app_data"))
	use_single_runner_invocation, _ := strconv.ParseBool(os.Getenv("use_single_runner_invocation"))
	use_mock_server, _ := strconv.ParseBool(os.Getenv("use_mock_server"))
	test_filters := os.Getenv("filter_test")

	// log.Print(strings.Split(os.Getenv("devices_list"), ","))

	payload := BrowserStackPayload{
		// Devices:                []string{os.Getenv("devices_list")},
		InstrumentationLogs:    instrumentation_logs,
		NetworkLogs:            network_logs,
		DeviceLogs:             device_logs,
		Screenshots:            screenshots,
		VideoRecording:         video_recording,
		SingleRunnerInvocation: use_single_runner_invocation,
		Project:                os.Getenv("project"),
		ProjectNotifyURL:       os.Getenv("project_notify_url"),
		UseLocal:               use_local,
		ClearAppData:           clear_app_data,
		UseMockServer:          use_mock_server,
		// Class:                  []string{os.Getenv("filter_test")},
		// Package:                []string{os.Getenv("filter_test")},
		// Annotation:             []string{os.Getenv("filter_test")},
		// Size:                   []string{os.Getenv("filter_test")},
	}

	scanner := bufio.NewScanner(strings.NewReader(os.Getenv("devices_list")))
	for scanner.Scan() {
		device := scanner.Text()
		device = strings.TrimSpace(device)

		if device == "" {
			continue
		}

		payload.Devices = append(payload.Devices, device)

	}

	if test_filters != "" {
		payload.Class = []string{test_filters}
		payload.Package = []string{test_filters}
	}

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

// var payload_validations map[string]interface{}

// // WIP
// func payload_validate_key(key string, value string) bool {
// 	fmt.Println(payload_validations[key])
// 	if validation, key_exists := payload_validations[key]; key_exists {
// 		// log.Print(validation, key_exists)

// 		validation_regex := regexp.MustCompile(validation.(string))
// 		return validation_regex.MatchString(value)
// 	}

// 	failf("Failed to parse the input parameter: %s", key)
// 	return false
// }

// // func payload_validate()

// func utils() {
// 	jsonFile, err := os.Open("validation.json")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer jsonFile.Close()

// 	byteValue, _ := ioutil.ReadAll(jsonFile)

// 	json.Unmarshal([]byte(byteValue), &payload_validations)
// }
