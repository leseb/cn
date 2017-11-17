package cmd

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
)

func DockerExist() {
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

func Selinux() {
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
