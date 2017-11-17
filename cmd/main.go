package cmd

var Version = "1.0.0"
var WorkingDirectory = "/usr/share/ceph-nano"
var CephNanoUid = "nano"
var RgwPort = "8000"
var CephNanoAccessKey = "accesskey"
var CephNanoSecretKey = "secretkey"
var ContainerName = "ceph-nano"

func Main() {
	DockerExist()
	NanoCli()
}
