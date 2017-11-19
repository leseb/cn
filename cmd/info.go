package cmd

import (
	"fmt"
)

// echoInfo prints useful information about Ceph Nano
func echoInfo() {
	// Always wait the container to be ready
	cephNanoHealth()
	cephNanoS3Health()

	// Fetch Amazon Keys
	getAwsKey()

	// Get IPs
	ips, _ := getInterfaceIPv4s()

	InfoLine := "\n" +
		"Ceph status is: $(docker exec ceph-nano ceph health) \n" +
		"Ceph Rados Gateway address is: http://" + ips[0].String() + ":8000 \n" +
		"Your working directory is: " +
		WorkingDirectory +
		"\n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		CephNanoSecretKey +
		"\n" +
		""
	fmt.Println(InfoLine)
}
