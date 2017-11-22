package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
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
// s3cmdWrapper wraps s3cmd commands to cn
func s3cmdWrapperMb(c *cli.Context) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}

	if len(c.Args()) > 1 {
		fmt.Printf("Too many arguments!")
		os.Exit(1)
	}
	var argFirst string
	argFirst = c.Args().First()

	cmd := []string{"s3cmd", "mb", "s3://" + argFirst}
	output := execContainer(ContainerName, cmd)
	fmt.Printf("%s", output)
}

func s3cmdWrapper(c *cli.Context) {
	fmt.Printf("%#v\n", c.Args().Tail())
}
