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
)

const browserstack_domain = "https://api-cloud.browserstack.com"
const app_upload_endpoint = "/app-automate/upload"
const test_suite_upload_endpoint = "/app-automate/espresso/v2/test-suite"
const app_automate_build_endpoint = "/app-automate/espresso/v2/build"

// func failf(format string, args ...interface{}) {
// 	log.Fatalf(format, args...)
// 	os.Exit(1)
// }

func main() {
	log.Printf("Hello from go")

	// todo: get from env
	username := "rishabhbhatia_OZ2u1M" // get from env
	access_key := "e76ypTPaVtQnFyqhAWBn"

	// upload the app
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
	upload_test_suite := upload("/Users/rishabh/Downloads/CalculatorTest.apk", test_suite_upload_endpoint, username, access_key)

	test_suite_parsed_response := make(map[string]interface{})
	test_suite_parse_err := json.Unmarshal([]byte(upload_test_suite), &test_suite_parsed_response)

	if test_suite_parse_err != nil {
		failf("Unable to parse test_suite_upload API response: %s", err)
	}

	test_suite_url := test_suite_parsed_response["test_suite_url"].(string)
	log.Println("test_suite_url -> ", test_suite_url)

	build(app_url, test_suite_url, username, access_key)
}

func upload(app_path string, endpoint string, username string, access_key string) string {
	if app_path == "" {
		failf("Failed to upload app on BrowserStack, error : app_path not found")
	}

	method := "POST"

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
	req, err := http.NewRequest(method, browserstack_domain+endpoint, payload)

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

func build(app_url string, test_suite_url string, username string, access_key string) string {
	if app_url == "" || test_suite_url == "" {
		failf("Failed to upload app on BrowserStack, error : app_path not found")
	}

	log.Printf("Starting build with app_id %s and test_suite_id", app_url, test_suite_url)

	method := "POST"

	devices := [...]string{"Samsung Galaxy Note 20-10.0"}

	payload_values := map[string]interface{}{
		"app":       app_url,
		"testSuite": test_suite_url,
		"devices":   devices,
	}

	payload, err := json.Marshal(payload_values)

	client := &http.Client{}
	req, err := http.NewRequest(method, browserstack_domain+app_automate_build_endpoint, bytes.NewBuffer(payload))

	if err != nil {
		// TODO: confirm this error
		failf("Failed to upload file: %s", err)
	}
	// req.SetBasicAuth("rishabhbhatia_OZ2u1M", "e76ypTPaVtQnFyqhAWBn")
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

	log.Print(string(body))
	return string(body)
}
