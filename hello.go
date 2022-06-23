package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func main() {
	log.Printf("Hello from go")

	// log.Print(os.Environ())

	username := os.Getenv("browserstack_username")

	access_key := os.Getenv("browserstack_accesskey")

	android_app := os.Getenv("android_app_under_test")
	test_suite := os.Getenv("espresso_test_suite")

	if username == "" || access_key == "" {
		failf("Failed to upload app on BrowserStack, error : invalid credentials")
	}

	// upload the app
	upload_app := upload(android_app, app_upload_endpoint, username, access_key)

	app_upload_parsed_response := make(map[string]interface{})

	// json decode
	err := json.Unmarshal([]byte(upload_app), &app_upload_parsed_response)

	if err != nil {
		failf("Unable to parse app_upload API response: %s", err)
	}

	app_url := app_upload_parsed_response["app_url"].(string)
	log.Println("app_url ->", app_url)

	upload_test_suite := upload(test_suite, test_suite_upload_endpoint, username, access_key)

	test_suite_parsed_response := make(map[string]interface{})
	test_suite_parse_err := json.Unmarshal([]byte(upload_test_suite), &test_suite_parsed_response)

	if test_suite_parse_err != nil {
		failf("Unable to parse test_suite_upload API response: %s", err)
	}

	test_suite_url := test_suite_parsed_response["test_suite_url"].(string)
	log.Println("test_suite_url -> ", test_suite_url)

	build_response := build(app_url, test_suite_url, username, access_key)

	check_build_status, _ := strconv.ParseBool(os.Getenv("check_build_status"))

	if check_build_status {
		build_parsed_response := make(map[string]interface{})

		json.Unmarshal([]byte(build_response), &build_parsed_response)

		build_id := build_parsed_response["build_id"].(string)

		build_status := checkBuildStatus(build_id, username, access_key)

		log.Print(build_status)
		return
	} else {
		return
	}

}
