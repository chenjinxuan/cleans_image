package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/gorhill/cronexpr"
	"os"
	"os/signal"
	"syscall"
	"time"
)


var (
	cronStr = flag.String("cron", "0 0 3 * MON *", "Perform the cleaning task regularly and follow the cron expression default 0 0 3 * MON *")
	once    = flag.Bool("once", false, "Perform the cleanup task only once")
	force   = flag.Bool("force", true, "Force image removal")
	tryRun  = flag.Bool("tryRun", false, "Try to run without deleting the image")
)

func main() {
	flag.Parse()
	err := cleans()
	if err != nil {
		fmt.Println("cleans fail " + err.Error())
	}
	if *once {
		return
	}
	expr, err := cronexpr.Parse(*cronStr)
	if err != nil {
		panic(err.Error())
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		now := time.Now()
		next := expr.Next(now)
		diff := next.Sub(now)
		select {
		case <-time.After(diff):
			err = cleans()
			if err != nil {
				fmt.Println("cleans fail " + err.Error())
			}
		case <-signalChan:
			return
		}
	}
}

func getClient() (*docker.Client, error) {
	if os.Getenv("DOCKER_API_VERSION") == "" {
		os.Setenv("DOCKER_API_VERSION", "1.23")
	}
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func cleans() error {
	cli, err := getClient()
	if err != nil {
		fmt.Println("get client fail " + err.Error())
		return err
	}
	cMap, err := getAllContainers(cli)
	if err != nil {
		fmt.Println("get container list fail " + err.Error())
		return err
	}
	images, err := getAllImages(cli)
	if err != nil {
		fmt.Println("get images list fail " + err.Error())
		return err
	}
	for _, image := range images {
		if _, ok := cMap[image.ID]; !ok {
			if !*tryRun {
				err = removeImages(cli, image.ID)
				if err != nil {
					fmt.Println(fmt.Sprintf("remove image ID %s fail %s", image.ID, err.Error()))
					continue
				}
				fmt.Println(fmt.Sprintf("remove image ID %s success", image.ID))
			} else {
				fmt.Println(fmt.Sprintf("need delete images ID %s", image.ID))
			}

		}
	}
	return nil
}

func getAllContainers(cli *docker.Client) (map[string]struct{}, error) {
	containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{All: false})
	if err != nil {
		return nil, err
	}
	cMap := make(map[string]struct{})
	for _, con := range containers {
		cMap[con.ImageID] = struct{}{}
	}
	return cMap, nil
}

func getAllImages(cli *docker.Client) ([]types.ImageSummary, error) {
	return cli.ImageList(context.TODO(), types.ImageListOptions{All: false})
}

func removeImages(cli *docker.Client, ID string) error {
	_, err := cli.ImageRemove(context.TODO(), ID, types.ImageRemoveOptions{Force: *force})
	return err
}
