// Package client exposes a Docker client wrapper for managing containers.
package client

import (
	"context"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// Container represents a Docker container with essential details.
type Container struct {
	container.Config
	ID    string `json:"Id"`
	Name  string `json:"Name"`
	Image string `json:"Image"`
	State string `json:"State"`
}

// ClientWrapper wraps the Docker client to provide container management functionalities.
type ClientWrapper struct {
	client *client.Client
}

// NewClient creates a new ClientWrapper with an initialized Docker client.
func NewClient() *ClientWrapper {
	dockerClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err.Error())
	}

	return &ClientWrapper{client: dockerClient}
}

// CloseClient closes the Docker client connection.
func (cw *ClientWrapper) CloseClient() {
	cw.client.Close()
}

// GetContainers retrieves a list of all Docker containers.
func (cw *ClientWrapper) GetContainers() []Container {
	containers, err := cw.client.ContainerList(
		context.Background(),
		client.ContainerListOptions{All: true},
	)
	if err != nil {
		panic(err)
	}
	var dockerContainers []Container
	for _, container := range containers.Items {
		dockerContainers = append(dockerContainers, Container{
			ID:    container.ID,
			Name:  container.Names[0][1:],
			Image: container.Image,
			State: string(container.State),
		})
	}
	return dockerContainers
}

// GetContainerState retrieves the current state of a specific Docker container by its ID.
func (cw *ClientWrapper) GetContainerState(id string) string {
	inspectResponse, _ := cw.client.ContainerInspect(context.Background(), id, client.ContainerInspectOptions{})
	return string(inspectResponse.Container.State.Status)
}

// PauseContainer pauses a specific Docker container by its ID.
func (cw *ClientWrapper) PauseContainer(id string) {
	cw.client.ContainerPause(context.Background(), id, client.ContainerPauseOptions{})
}

// PauseContainers pauses multiple Docker containers by their IDs.
func (cw *ClientWrapper) PauseContainers(ids []string) {
	for _, id := range ids {
		cw.PauseContainer(id)
	}
}

// UnpauseContainer unpauses a specific Docker container by its ID.
func (cw *ClientWrapper) UnpauseContainer(id string) {
	cw.client.ContainerUnpause(context.Background(), id, client.ContainerUnpauseOptions{})
}

// UnpauseContainers unpauses multiple Docker containers by their IDs.
func (cw *ClientWrapper) UnpauseContainers(ids []string) {
	for _, id := range ids {
		cw.UnpauseContainer(id)
	}
}

// StartContainer starts a specific Docker container by its ID.
func (cw *ClientWrapper) StartContainer(id string) {
	cw.client.ContainerStart(context.Background(), id, client.ContainerStartOptions{})
}

// StartContainers starts multiple Docker containers by their IDs.
func (cw *ClientWrapper) StartContainers(ids []string) {
	for _, id := range ids {
		cw.StartContainer(id)
	}
}

// StopContainer stops a specific Docker container by its ID.
func (cw *ClientWrapper) StopContainer(id string) {
	cw.client.ContainerStop(context.Background(), id, client.ContainerStopOptions{})
}

// StopContainers stops multiple Docker containers by their IDs.
func (cw *ClientWrapper) StopContainers(ids []string) {
	for _, id := range ids {
		cw.StopContainer(id)
	}
}
