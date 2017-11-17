package cmd

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"strings"
)

func GrepForSuccess() bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	out, err := cli.ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()

	if strings.Contains(newStr, "SUCCESS") {
		return true
	}
	return false
}

func CephNanoHealth() {
	// remove that once the loop works
	GrepForSuccess()

	/*
		timeout := time.After(30 * time.Second)
		poll := 3 * time.Second
			for {
				Health := GrepForSuccess()
				select {
				case <-Health:
					fmt.Println("The end!")
					return
				case <-timeout:
					fmt.Println("There's no more time to this. Exiting!")
					return
				default:
					fmt.Println("still waiting")
				}
				time.Sleep(poll)
			}
	*/
}

func CephNanoS3Health() {
	// curl --fail --silent --output /dev/null http://"$IP":8000
}

func WaitForContainer() {
	CephNanoHealth()
	CephNanoS3Health()
}
