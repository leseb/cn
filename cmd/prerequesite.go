package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// DockerExist makes sure Docker is installed
func dockerExist() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println("Docker is not present on your system or not started.\n" +
			"Make sure it's started or follow installation instructions at https://docs.docker.com/engine/installation/")
		os.Exit(1)
	}
}

// Selinux checks if Selinux is installed and set to Enforcing,
// we relabel our WorkingDirectory to allow the container to access files in this directory
func selinux() {
	if _, err := os.Stat("/sbin/getenforce"); !os.IsNotExist(err) {
		out, err := exec.Command("getenforce").Output()
		if err != nil {
			log.Fatal(err)
		}
		if string(out) == "Enforcing" {
			exec.Command("chcon -Rt svirt_sandbox_file_t %s", WorkingDirectory)
		}
	}
}
