package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func failf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
	os.Exit(1)
}

var payload_validations map[string]interface{}

// WIP
func payload_validate_key(key string, value string) bool {
	fmt.Println(payload_validations[key])
	if validation, key_exists := payload_validations[key]; key_exists {
		// log.Print(validation, key_exists)

		validation_regex := regexp.MustCompile(validation.(string))
		return validation_regex.MatchString(value)
	}

	failf("Failed to parse the input parameter: %s", key)
	return false
}

// func payload_validate()

func main() {
	jsonFile, err := os.Open("validation.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &payload_validations)
}
