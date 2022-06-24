package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	// 	fmt.Println("This is the value specified for the input 'example_step_input':", os.Getenv("example_step_input"))

	// 	//
	// 	// --- Step Outputs: Export Environment Variables for other Steps:
	// 	// You can export Environment Variables for other Steps with
	// 	//  envman, which is automatically installed by `bitrise setup`.
	// 	// A very simple example:
	// 	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "EXAMPLE_STEP_OUTPUT", "--value", "the value you want to share").CombinedOutput()
	// 	if err != nil {
	// 		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
	// 		os.Exit(1)
	// 	}
	// 	// You can find more usage examples on envman's GitHub page
	// 	//  at: https://github.com/bitrise-io/envman

	// 	//
	// 	// --- Exit codes:
	// 	// The exit code of your Step is very important. If you return
	// 	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// 	// Any non zero exit code will be registered as "failed" by `bitrise`.
	// 	os.Exit(0)

	username := os.Getenv("browserstack_username")
	access_key := os.Getenv("browserstack_accesskey")
	android_app := os.Getenv("android_app_under_test")
	test_suite := os.Getenv("espresso_test_suite")

	if username == "" || access_key == "" {
		failf(UPLOAD_APP_ERROR, "invalid credentials")
	}

	upload_app, err := upload(android_app, APP_UPLOAD_ENDPOINT, username, access_key)

	if err != nil {
		failf(err.Error())
	}

	upload_app_parsed_response := jsonParse(upload_app)

	if upload_app_parsed_response["app_url"] == "" {
		failf(err.Error())
	}

	app_url := upload_app_parsed_response["app_url"].(string)

	upload_test_suite, err := upload(test_suite, TEST_SUITE_UPLOAD_ENDPOINT, username, access_key)

	if err != nil {
		failf(err.Error())
	}

	test_suite_url := jsonParse(upload_test_suite)["test_suite_url"].(string)

	build_response, err := build(app_url, test_suite_url, username, access_key)

	if err != nil {
		failf(err.Error())
	}

	check_build_status, _ := strconv.ParseBool(os.Getenv("check_build_status"))

	if check_build_status {
		build_parsed_response := jsonParse(build_response)

		if build_parsed_response["message"] != "Success" {
			failf(BUILD_FAILED_ERROR, build_parsed_response["message"])
		}

		build_id := build_parsed_response["build_id"].(string)

		build_status, err := checkBuildStatus(build_id, username, access_key)

		if err != nil {
			failf(err.Error())
		}

		log.Printf("build_status %s", build_status)
	}

	os.Exit(0)
}
