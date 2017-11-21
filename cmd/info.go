package cmd

import (
	"fmt"
	"strings"
)

// echoInfo prints useful information about Ceph Nano
func echoInfo() {
	// Always wait the container to be ready
	cephNanoHealth()
	cephNanoS3Health()

	// Fetch Amazon Keys
	//CephNanoAccessKey, CephNanoSecretKey := getAwsKey()

	// Get Ceph health
	cmd := []string{"ceph", "health"}
	c := execContainer(ContainerName, cmd)

	// Get IPs
	ips, _ := getInterfaceIPv4s()

	InfoLine := strings.TrimSpace(string(c)) +
		" is the Ceph status. \n" +
		"Ceph Rados Gateway address is: http://" + ips[0].String() + ":8000 \n" +
		"Your working directory is: " +
		WorkingDirectory +
		"\n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		//CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		//CephNanoSecretKey +
		"\n" +
		""
	fmt.Println(InfoLine)
}
