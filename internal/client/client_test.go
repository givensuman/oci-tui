package client

import (
	"context"
	"io"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/google/uuid"
)

func TestNewClient(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Client.Ping(context.Background())
	if err != nil {
		fmt.Println("failed to ping daemon")
		t.Error(err)
	}
}

func TestContainerList(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		fmt.Println("failed to call ContainerList")
		t.Error(err)
	}
}

func TestContainerCreateListAndRemove(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	containerName, err := uuid.NewUUID()
	if err != nil {
		t.FailNow()
	}

	imageName := "hello-world"
	reader, err := cli.Client.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		fmt.Printf("failed to pull image %s\n", imageName)
		t.Skip(err)
	}

	io.Copy(io.Discard, reader)
	config := &container.Config{
		Image: imageName,
	}

	createResp, err := cli.Client.ContainerCreate(context.Background(), config, nil, nil, nil, containerName.String())
	if err != nil {
		t.Error(err)
	}

	listResp, err := cli.Client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, c := range listResp {
		if c.ID == createResp.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created container %s not found in list", createResp.ID)
	}

	cli.Client.ContainerRemove(context.Background(), createResp.ID, container.RemoveOptions{})
}
