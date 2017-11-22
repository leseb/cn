package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdSync is the Cobra CLI call
func CliS3CmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync LOCAL_DIR BUCKET[/PREFIX]",
		Short: "Synchronize a directory tree to S3",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdSync,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdSync wraps s3cmd command in the container
func S3CmdSync(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "sync", args[0], "s3://" + args[1]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
