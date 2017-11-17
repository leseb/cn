package cmd

import (
	"fmt"
)

func EchoInfo() {
	// Always wait the container to be ready
	WaitForContainer()

	// Fetch Amazon Keys
	GetAwsKey()

	InfoLine := "\n" +
		"Ceph status is: $(docker exec ceph-nano ceph health) \n" +
		"Ceph Rados Gateway address is: http://$IP:$RGW_CIVETWEB_PORT \n" +
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
