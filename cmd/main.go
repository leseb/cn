package cmd

// Version is the Ceph Nano version
var Version = "1.0.0"

// WorkingDirectory is the working directory where objects can be put inside S3
var WorkingDirectory = "/usr/share/ceph-nano"

// CephNanoUID is the uid of the S3 user
var CephNanoUID = "nano"

// RgwPort is the rgw listenning port
var RgwPort = "8000"

// ContainerName is name of the container
var ContainerName = "ceph-nano"

// Main is the main function calling the whole program
func Main() {
	dockerExist()
	nanoCli()
}
