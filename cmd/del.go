package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
function s3cmd_wrap {
  if docker ps | grep -sq ceph-nano; then
    IFS=" " read -r -a array <<< "$*"
    local docker_exec
    docker_exec="docker exec ceph-nano s3cmd"
    if [[ "${array[0]}" =~ mb|rb|ls|del|info|du ]]; then
      $docker_exec "${array[0]}" "${array[1]/#/s3://}"
    elif [[ "${array[0]}" =~ cp|mv ]]; then
      $docker_exec "${array[0]}" "${array[1]/#/s3://}" "${array[2]/#/s3://}"
    elif [[ "${array[0]}" =~ get ]]; then
      $docker_exec "${array[0]}" "${array[1]/#/s3://}" "${array[2]}"
    elif [[ "${array[0]}" =~ put|sync ]]; then
      local input_file="${array[1]}"
      if [ ! -e "$input_file" ]; then
          echo "$input_file doesn't exist !"
          return
      fi
      if [ "$(dirname $input_file)" != "$WORKING_DIR" ]; then
          echo "$input_file should be in $WORKING_DIR directory !"
          return
      fi
      $docker_exec "${array[0]}" "${array[1]}" "${array[2]/#/s3://}"
    else
      $docker_exec "${array[@]}"
    fi
  else
    echo "$PROGRAM is not running so S3 calls are not avaiable."
    echo "Start it with: ./$PROGRAM start [working dir]"
  fi
}
*/

// CliS3CmdDel is the Cobra CLI call
func CliS3CmdDel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del [BUCKET]/OBJECT",
		Short: "Make bucket",
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
