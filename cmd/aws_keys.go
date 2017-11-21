package cmd

import (
	"encoding/json"
)

// getAwsKey gets AWS keys from inside the container
func getAwsKey() (string, string) {
	cmd := []string{"/bin/cat", "/nano_user_details"}
	output := execContainer(ContainerName, cmd)

	// declare structures for json
	type s3Details []struct {
		AccessKey string `json:"Access_key"`
		SecretKey string `json:"Secret_key"`
	}
	type jason struct {
		Keys s3Details
	}
	// assign variable to our json struct
	var parsedMap jason

	json.Unmarshal(output, &parsedMap)

	var CephNanoAccessKey string
	CephNanoAccessKey = parsedMap.Keys[0].AccessKey
	var CephNanoSecretKey string
	CephNanoSecretKey = parsedMap.Keys[0].SecretKey
	return CephNanoAccessKey, CephNanoSecretKey
}
