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
	CephNanoAccessKey, CephNanoSecretKey := getAwsKey()

	// Get Ceph health
	cmd := []string{"ceph", "health"}
	c := execContainer(ContainerName, cmd)

	// Get IPs, later using the first IP of the list is not ideal
	// However, Docker binds RGW port on 0.0.0.0 so any address will work
	ips, _ := getInterfaceIPv4s()

	InfoLine := "\n" +
		strings.TrimSpace(string(c)) +
		" is the Ceph status. \n" +
		"S3 server address is: http://" + ips[0].String() + ":8000 \n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		CephNanoSecretKey +
		"\n" +
		"Your working directory is: " +
		WorkingDirectory +
		"\n" +
		""
	fmt.Println(InfoLine)
}
