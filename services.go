package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func build(app_url string, test_suite_url string, username string, access_key string) (string, error) {
	if app_url == "" || test_suite_url == "" {
		return "", errors.New(FILE_NOT_AVAILABLE_ERROR)
	}

	payload_values := createBuildPayload()
	payload_values.App = app_url
	payload_values.TestSuite = test_suite_url

	payload, err := json.Marshal(payload_values)
	log.Print("Payload -> ", string(payload))

	final_payload := appendExtraCapabilities(string(payload))

	log.Print("Final payload -> ", string(final_payload))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", BROWSERSTACK_DOMAIN+APP_AUTOMATE_BUILD_ENDPOINT, bytes.NewBuffer(payload))

	req.SetBasicAuth(username+"-bitrise", access_key)

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		// Todo: confirm this error
		return "", errors.New(fmt.Sprintf(BUILD_FAILED_ERROR, err))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// Todo: confirm this error
		return "", errors.New(fmt.Sprintf(BUILD_FAILED_ERROR, err))
	}

	log.Print(string(body))

	return string(body), nil
}

// this function uploads both app and test suite
func upload(app_path string, endpoint string, username string, access_key string) (string, error) {
	if app_path == "" {
		return "", errors.New(FILE_NOT_AVAILABLE_ERROR)
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
		return "", errors.New(FILE_NOT_AVAILABLE_ERROR)
	}

	err := multipart_writer.Close()

	if err != nil {
		return "", errors.New(INVALID_FILE_TYPE_ERROR)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", BROWSERSTACK_DOMAIN+endpoint, payload)

	req.SetBasicAuth(username, access_key)

	req.Header.Set("Content-Type", multipart_writer.FormDataContentType())

	res, err := client.Do(req)

	if err != nil {
		return "", errors.New(fmt.Sprintf(UPLOAD_APP_ERROR, err))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", errors.New(fmt.Sprintf(UPLOAD_APP_ERROR, err))
	}

	return string(body), nil
}

func checkBuildStatus(build_id string, username string, access_key string) (string, error) {

	if build_id == "" {
		return "", errors.New(fmt.Sprintf(FETCH_BUILD_STATUS_ERROR, "invalid build_id"))
	}

	build_parsed_response := make(map[string]interface{})
	build_status := ""
	var body []byte
	var build_status_error error

	clear := setInterval(func() {

		client := &http.Client{}
		req, _ := http.NewRequest("GET", BROWSERSTACK_DOMAIN+APP_AUTOMATE_BUILD_STATUS_ENDPOINT+build_id, nil)

		req.SetBasicAuth(username, access_key)

		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)

		if err != nil {
			build_status_error = errors.New(fmt.Sprintf(FETCH_BUILD_STATUS_ERROR, err))
			return
		}

		defer res.Body.Close()

		body, err = ioutil.ReadAll(res.Body)

		if err != nil {
			build_status_error = errors.New(fmt.Sprintf(FETCH_BUILD_STATUS_ERROR, err))
			return
		}

		json.Unmarshal([]byte(body), &build_parsed_response)

		if build_parsed_response["error"] != nil && build_parsed_response["error"] != "" {
			build_status_error = errors.New(fmt.Sprintf(FETCH_BUILD_STATUS_ERROR, build_parsed_response["error"]))
			return
		}

		build_status = build_parsed_response["status"].(string)

	}, 10, false)

	for {
		if build_status != "running" && build_status != "" {
			// Stop the ticket, ending the interval go routine
			clear <- true
			printBuildStatus(build_parsed_response)

			return build_status, build_status_error
		}

		if build_status_error != nil {

			clear <- true

			return build_status, build_status_error
		}
	}

}
