package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdGet is the Cobra CLI call
func CliS3CmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get BUCKET/OBJECT LOCAL_FILE",
		Short: "Get file into bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdGet,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdGet wraps s3cmd command in the container
func S3CmdGet(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "get", "s3://" + args[0], args[1]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
