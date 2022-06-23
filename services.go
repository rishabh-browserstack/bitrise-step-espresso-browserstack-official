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
const app_automate_build_status_endpoint = "/app-automate/espresso/v2/builds/"

const interval_in_milliseconds = 30000

func build(app_url string, test_suite_url string, username string, access_key string) string {
	if app_url == "" || test_suite_url == "" {
		failf("Failed to upload app on BrowserStack, error : app_path not found")
	}

	payload_values := createBuildPayload()
	payload_values.App = app_url
	payload_values.TestSuite = test_suite_url

	log.Print(payload_values)

	payload, err := json.Marshal(payload_values)

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
	}, interval_in_milliseconds, false)

	log.Print(build_status)

	for {
		if build_status != "running" && build_status != "" {
			// Stop the ticket, ending the interval go routine
			clear <- true
			return build_status
		}
	}
}
