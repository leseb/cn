package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// logsNano prints rgw logs
func logsNano(c *cli.Context) {
	showS3Logs()
}

func showS3Logs() {
	c := []string{"cat", "/var/log/ceph/client.rgw.ceph-nano-faa32aebf00b.log"}
	output := execContainer(ContainerName, c)
	fmt.Printf("%s", output)
}
