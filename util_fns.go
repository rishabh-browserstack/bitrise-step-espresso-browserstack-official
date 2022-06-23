package main

import (
	"log"
	"os"
	"strconv"
)

// type Queue []Class

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

	payload := BrowserStackPayload{
		Devices:                []string{os.Getenv("devices_list")},
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
