package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const browserstack_domain = "https://api-cloud.browserstack.com"
const app_upload_endpoint = "/app-automate/upload"
const test_suite_upload_endpoint = "/app-automate/espresso/v2/test-suite"
const app_automate_build_endpoint = "/app-automate/espresso/v2/build"
const app_automate_build_status_endpoint = "/app-automate/espresso/v2/builds/"

func build(app_url string, test_suite_url string, username string, access_key string) string {
	if app_url == "" || test_suite_url == "" {
		failf("Failed to upload app on BrowserStack, error : app_path not found")
	}

	log.Printf("Starting build with app_id %s and test_suite_id", app_url, test_suite_url)

	// devices := [...]string{"Samsung Galaxy Note 20-10.0"}

	// payload_values := map[string]interface{}{
	// 	"app":       app_url,
	// 	"testSuite": test_suite_url,
	// 	"devices":   []string{os.Getenv("devices_list")},
	// 	"deviceLogs":
	// }

	payload_values := createBuildPayload()
	payload_values.App = app_url
	payload_values.TestSuite = test_suite_url

	// log.Print(os.Getenv("devices_list"))
	log.Print(payload_values)

	payload, err := json.MarshalIndent(payload_values, "", "  ")

	client := &http.Client{}
	req, err := http.NewRequest("POST", browserstack_domain+app_automate_build_endpoint, bytes.NewBuffer(payload))

	if err != nil {
		// TODO: confirm this error
		failf("Failed to upload file: %s", err)
	}

	req.SetBasicAuth(username+"-bitrise", access_key)

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		failf("Unable to read response: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		failf("Unable to read api response: %s", err)
	}

	log.Print(string(body))

	return string(body)
}

func main() {
	log.Printf("Hello from go")

	// log.Print(os.Environ())

	username := os.Getenv("browserstack_username")
	access_key := os.Getenv("browserstack_accesskey")

	if username == "" || access_key == "" {
		failf("Failed to upload app on BrowserStack, error : invalid credentials")
	}

	// upload the app
	// todo: get from env of bitrise
	upload_app := upload("/Users/rishabh/Downloads/Calculator.apk", app_upload_endpoint, username, access_key)

	app_upload_parsed_response := make(map[string]interface{})

	// json decode
	err := json.Unmarshal([]byte(upload_app), &app_upload_parsed_response)

	if err != nil {
		failf("Unable to parse app_upload API response: %s", err)
	}

	app_url := app_upload_parsed_response["app_url"].(string)
	log.Println("app_url ->", app_url)

	// upload test_suite
	// todo: get from env of bitrise
	upload_test_suite := upload("/Users/rishabh/Downloads/CalculatorTest.apk", test_suite_upload_endpoint, username, access_key)

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
	} else {
		return
	}

}

func upload(app_path string, endpoint string, username string, access_key string) string {
	if app_path == "" {
		failf("Failed to upload app on BrowserStack, error : app_path not found")
	}

	payload := &bytes.Buffer{}
	multipart_writer := multipart.NewWriter(payload)
	file, fileErr := os.Open(app_path)

	defer file.Close()

	// creates a new form data header
	// reading and copying the file's content to form data
	attached_file,
		fileErr := multipart_writer.CreateFormFile("file", filepath.Base(app_path))

	_, fileErr = io.Copy(attached_file, file)

	if fileErr != nil {
		// TODO: confirm this error
		failf("Unable to read file: %s", fileErr)
	}

	err := multipart_writer.Close()
	if err != nil {
		// TODO: confirm this error
		failf("Unable to close file: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", browserstack_domain+endpoint, payload)

	if err != nil {
		// TODO: confirm this error
		failf("Failed to upload file: %s", err)
	}

	req.SetBasicAuth(username, access_key)

	req.Header.Set("Content-Type", multipart_writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		failf("Unable to read response: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		failf("Unable to read api response: %s", err)
	}

	return string(body)
}

func checkBuildStatus(build_id string, username string, access_key string) string {

	build_status := ""

	clear := setInterval(func() {
		log.Print("Inside interval function")
		log.Printf("Checking build status for build_id %s", build_id)

		client := &http.Client{}
		req, err := http.NewRequest("GET", browserstack_domain+app_automate_build_status_endpoint+build_id, nil)

		log.Print(browserstack_domain + app_automate_build_status_endpoint + build_id)
		if err != nil {
			// TODO: confirm this error
			failf("Failed to check build status: %s", err)
		}

		req.SetBasicAuth(username, access_key)

		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)

		if err != nil {
			failf("Unable to read response: %s", err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			failf("Unable to read api response: %s", err)
		}
		// return string(body)

		build_parsed_response := make(map[string]interface{})

		json.Unmarshal([]byte(body), &build_parsed_response)

		build_status = build_parsed_response["status"].(string)
		log.Print(build_status)
	}, 30000, false)

	log.Print(build_status)

	for {
		if build_status != "running" && build_status != "" {
			// Stop the ticket, ending the interval go routine
			clear <- true
			return build_status
		}
	}
}
