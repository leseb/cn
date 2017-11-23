package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdLs is the Cobra CLI call
func CliS3CmdLs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls [BUCKET]",
		Short: "List objects or buckets",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdLs,
	}
	return cmd
}

// S3CmdLs wraps s3cmd command in the container
func S3CmdLs(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "ls", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
