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
		Args:  cobra.RangeArgs(2, 3),
		Run:   S3CmdGet,
		DisableFlagsInUseLine: true,
	}
	var Force string
	var Skip string
	var Continue string
	cmd.Flags().StringVarP(&Force, "force", "f", "true", "Force the get if the file already exists")
	cmd.Flags().StringVarP(&Skip, "skip", "s", "true", "Skip the get if the file already exists")
	cmd.Flags().StringVarP(&Continue, "continue", "c", "true", "Continue the get if the file already exists")

	return cmd
}

// S3CmdGet wraps s3cmd command in the container
func S3CmdGet(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}

	cmd.Help()

	// TODO REDO
	if len(args) > 2 {
		command := []string{"s3cmd", "get", "s3://" + args[0], "/tmp/" + args[1], args[2]}
		output := execContainer(ContainerName, command)
		fmt.Printf("%s", output)
	} else {
		command := []string{"s3cmd", "get", "s3://" + args[0], "/tmp/" + args[1]}
		output := execContainer(ContainerName, command)
		fmt.Printf("%s", output)
	}
}
