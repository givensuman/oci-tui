// Package client exposes a Docker client wrapper for managing containers.
package client

import (
	"context"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// Container represents a Docker container with essential details.
type Container struct {
	container.Config
	ID    string                   `json:"Id"`
	Name  string                   `json:"Name"`
	Image string                   `json:"Image"`
	State container.ContainerState `json:"State"`
}

// Image represents a Docker image.
type Image struct {
	ID       string   `json:"Id"`
	RepoTags []string `json:"RepoTags"`
	Size     int64    `json:"Size"`
	Created  int64    `json:"Created"`
}

// Network represents a Docker network.
type Network struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Driver string `json:"Driver"`
	Scope  string `json:"Scope"`
}

// Volume represents a Docker volume.
type Volume struct {
	Name       string `json:"Name"`
	Driver     string `json:"Driver"`
	Mountpoint string `json:"Mountpoint"`
}

// ClientWrapper wraps the Docker client to provide container management functionalities.
type ClientWrapper struct {
	client *client.Client
}

// NewClient creates a new ClientWrapper with an initialized Docker client.
func NewClient() (*ClientWrapper, error) {
	dockerClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &ClientWrapper{client: dockerClient}, nil
}

// CloseClient closes the Docker client connection.
func (cw *ClientWrapper) CloseClient() error {
	return cw.client.Close()
}

// GetContainers retrieves a list of all Docker containers.
func (cw *ClientWrapper) GetContainers() ([]Container, error) {
	listOptions := client.ContainerListOptions{
		All: true,
	}

	containers, err := cw.client.ContainerList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerContainers := make([]Container, 0, len(containers.Items))
	for _, container := range containers.Items {
		dockerContainers = append(dockerContainers, Container{
			ID:    container.ID,
			Name:  container.Names[0][1:],
			Image: container.Image,
			State: container.State,
		})
	}

	return dockerContainers, nil
}

// GetImages retrieves a list of all Docker images.
func (cw *ClientWrapper) GetImages() ([]Image, error) {
	// Assuming client.ImageListOptions exists and follows the pattern
	listOptions := client.ImageListOptions{
		All: true,
	}

	images, err := cw.client.ImageList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerImages := make([]Image, 0, len(images.Items))
	for _, img := range images.Items {
		dockerImages = append(dockerImages, Image{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
		})
	}

	return dockerImages, nil
}

// GetNetworks retrieves a list of all Docker networks.
func (cw *ClientWrapper) GetNetworks() ([]Network, error) {
	listOptions := client.NetworkListOptions{}

	networks, err := cw.client.NetworkList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerNetworks := make([]Network, 0, len(networks.Items))
	for _, net := range networks.Items {
		dockerNetworks = append(dockerNetworks, Network{
			ID:     net.ID,
			Name:   net.Name,
			Driver: net.Driver,
			Scope:  net.Scope,
		})
	}

	return dockerNetworks, nil
}

// GetVolumes retrieves a list of all Docker volumes.
func (cw *ClientWrapper) GetVolumes() ([]Volume, error) {
	listOptions := client.VolumeListOptions{}

	volumes, err := cw.client.VolumeList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerVolumes := make([]Volume, 0, len(volumes.Items))
	for _, vol := range volumes.Items {
		dockerVolumes = append(dockerVolumes, Volume{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
		})
	}

	return dockerVolumes, nil
}

// GetContainerState retrieves the current state of a specific Docker container by its ID.
func (cw *ClientWrapper) GetContainerState(id string) (string, error) {
	inspectResponse, err := cw.client.ContainerInspect(context.Background(), id, client.ContainerInspectOptions{})
	if err != nil {
		return "unknown", err
	}

	return string(inspectResponse.Container.State.Status), nil
}

// PauseContainer pauses a specific Docker container by its ID.
func (cw *ClientWrapper) PauseContainer(id string) error {
	_, err := cw.client.ContainerPause(context.Background(), id, client.ContainerPauseOptions{})
	return err
}

// PauseContainers pauses multiple Docker containers by their IDs.
func (cw *ClientWrapper) PauseContainers(ids []string) error {
	for _, id := range ids {
		if err := cw.PauseContainer(id); err != nil {
			return err
		}
	}

	return nil
}

// UnpauseContainer unpauses a specific Docker container by its ID.
func (cw *ClientWrapper) UnpauseContainer(id string) error {
	_, err := cw.client.ContainerUnpause(context.Background(), id, client.ContainerUnpauseOptions{})
	return err
}

// UnpauseContainers unpauses multiple Docker containers by their IDs.
func (cw *ClientWrapper) UnpauseContainers(ids []string) error {
	for _, id := range ids {
		if err := cw.UnpauseContainer(id); err != nil {
			return err
		}
	}

	return nil
}

// StartContainer starts a specific Docker container by its ID.
func (cw *ClientWrapper) StartContainer(id string) error {
	_, err := cw.client.ContainerStart(context.Background(), id, client.ContainerStartOptions{})
	return err
}

// StartContainers starts multiple Docker containers by their IDs.
func (cw *ClientWrapper) StartContainers(ids []string) error {
	for _, id := range ids {
		if err := cw.StartContainer(id); err != nil {
			return err
		}
	}

	return nil
}

// StopContainer stops a specific Docker container by its ID.
func (cw *ClientWrapper) StopContainer(id string) error {
	_, err := cw.client.ContainerStop(context.Background(), id, client.ContainerStopOptions{})
	return err
}

// StopContainers stops multiple Docker containers by their IDs.
func (cw *ClientWrapper) StopContainers(ids []string) error {
	for _, id := range ids {
		if err := cw.StopContainer(id); err != nil {
			return err
		}
	}

	return nil
}

// RemoveContainer removes a specific Docker container by its ID.
func (cw *ClientWrapper) RemoveContainer(id string) error {
	removeOptions := client.ContainerRemoveOptions{
		Force: true,
	}

	_, err := cw.client.ContainerRemove(context.Background(), id, removeOptions)
	return err
}

// RemoveContainers removes multiple Docker containers by their IDs.
func (cw *ClientWrapper) RemoveContainers(ids []string) error {
	for _, id := range ids {
		if err := cw.RemoveContainer(id); err != nil {
			return err
		}
	}

	return nil
}

// Logs represents the response from Moby's ContainerLogs.
type Logs client.ContainerLogsResult

// OpenLogs streams logs from a Docker container.
func (cw *ClientWrapper) OpenLogs(id string) (Logs, error) {
	logsOptions := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "all",
	}

	reader, err := cw.client.ContainerLogs(context.Background(), id, logsOptions)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// ExecShell starts an interactive shell (e.g., /bin/sh or /bin/bash) in the container with a TTY.
// Returns an io.ReadWriteCloser for bi-directional communication, or error.
func (cw *ClientWrapper) ExecShell(id string, shell []string) (io.ReadWriteCloser, error) {
	execCreateOptions := client.ExecCreateOptions{
		Cmd:          shell,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          true,
	}

	execResp, err := cw.client.ExecCreate(context.Background(), id, execCreateOptions)
	if err != nil {
		return nil, err
	}

	execAttachOptions := client.ExecAttachOptions{
		TTY: true,
	}

	attachResp, err := cw.client.ExecAttach(context.Background(), execResp.ID, execAttachOptions)
	if err != nil {
		return nil, err
	}

	return attachResp.Conn, nil // attaches to socket, full duplex
}
