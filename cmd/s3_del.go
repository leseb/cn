package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdDel is the Cobra CLI call
func CliS3CmdDel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del [BUCKET]/OBJECT",
		Short: "Delete bucket",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdDel,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdDel wraps s3cmd command in the container
func S3CmdDel(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "del", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
