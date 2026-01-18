// Package client exposes a Docker client wrapper for managing containers.
package client

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// ContainerStats represents the CPU and memory usage of a container.
type ContainerStats struct {
	CPUPercent float64
	MemUsage   float64
	MemLimit   float64
	NetRx      float64
	NetTx      float64
}

// Container represents a Docker container with essential details.
type Container struct {
	container.Config
	ID    string `json:"Id"`
	Name  string `json:"Name"`
	Image string `json:"Image"`
	State string `json:"State"`
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
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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
	listOptions := container.ListOptions{
		All: true,
	}

	containers, err := cw.client.ContainerList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerContainers := make([]Container, 0, len(containers))
	for _, c := range containers {
		dockerContainers = append(dockerContainers, Container{
			ID:    c.ID,
			Name:  c.Names[0][1:],
			Image: c.Image,
			State: c.State,
		})
	}

	return dockerContainers, nil
}

// GetImages retrieves a list of all Docker images.
func (cw *ClientWrapper) GetImages() ([]Image, error) {
	listOptions := types.ImageListOptions{
		All: true,
	}

	images, err := cw.client.ImageList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerImages := make([]Image, 0, len(images))
	for _, img := range images {
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
	listOptions := types.NetworkListOptions{}

	networks, err := cw.client.NetworkList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerNetworks := make([]Network, 0, len(networks))
	for _, net := range networks {
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
	listOptions := volume.ListOptions{}

	volumes, err := cw.client.VolumeList(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	dockerVolumes := make([]Volume, 0, len(volumes.Volumes))
	for _, vol := range volumes.Volumes {
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
	inspectResponse, err := cw.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return "unknown", err
	}

	return string(inspectResponse.State.Status), nil
}

// PauseContainer pauses a specific Docker container by its ID.
func (cw *ClientWrapper) PauseContainer(id string) error {
	return cw.client.ContainerPause(context.Background(), id)
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
	return cw.client.ContainerUnpause(context.Background(), id)
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
	return cw.client.ContainerStart(context.Background(), id, container.StartOptions{})
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
	return cw.client.ContainerStop(context.Background(), id, container.StopOptions{})
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
	removeOptions := container.RemoveOptions{
		Force: true,
	}

	return cw.client.ContainerRemove(context.Background(), id, removeOptions)
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
type Logs io.ReadCloser

// OpenLogs streams logs from a Docker container.
func (cw *ClientWrapper) OpenLogs(id string) (Logs, error) {
	logsOptions := container.LogsOptions{
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
	execCreateOptions := types.ExecConfig{
		Cmd:          shell,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execResp, err := cw.client.ContainerExecCreate(context.Background(), id, execCreateOptions)
	if err != nil {
		return nil, err
	}

	execAttachOptions := types.ExecStartCheck{
		Tty: true,
	}

	attachResp, err := cw.client.ContainerExecAttach(context.Background(), execResp.ID, execAttachOptions)
	if err != nil {
		return nil, err
	}

	return attachResp.Conn, nil // attaches to socket, full duplex
}

// RemoveImage removes a specific Docker image by its ID.
func (cw *ClientWrapper) RemoveImage(id string) error {
	options := types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: true,
	}

	_, err := cw.client.ImageRemove(context.Background(), id, options)
	return err
}

// RemoveVolume removes a specific Docker volume by its name.
func (cw *ClientWrapper) RemoveVolume(name string) error {
	return cw.client.VolumeRemove(context.Background(), name, false)
}

// RemoveNetwork removes a specific Docker network by its ID.
func (cw *ClientWrapper) RemoveNetwork(id string) error {
	return cw.client.NetworkRemove(context.Background(), id)
}

// PruneImages removes all unused images.
func (cw *ClientWrapper) PruneImages() (uint64, error) {
	report, err := cw.client.ImagesPrune(context.Background(), filters.Args{})
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// PruneVolumes removes all unused volumes.
func (cw *ClientWrapper) PruneVolumes() (uint64, error) {
	report, err := cw.client.VolumesPrune(context.Background(), filters.Args{})
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// PruneNetworks removes all unused networks.
func (cw *ClientWrapper) PruneNetworks() error {
	_, err := cw.client.NetworksPrune(context.Background(), filters.Args{})
	return err
}

// GetContainersUsingImage returns a list of container names that are using the specified image ID.
func (cw *ClientWrapper) GetContainersUsingImage(imageID string) ([]string, error) {
	containers, err := cw.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var usedBy []string
	for _, c := range containers {
		if c.ImageID == imageID {
			// Name usually comes with a slash, e.g., "/my-container"
			name := c.Names[0]
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
			usedBy = append(usedBy, name)
		}
	}
	return usedBy, nil
}

// GetContainersUsingVolume returns a list of container names that are using the specified volume name.
func (cw *ClientWrapper) GetContainersUsingVolume(volumeName string) ([]string, error) {
	containers, err := cw.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var usedBy []string
	for _, c := range containers {
		for _, m := range c.Mounts {
			if m.Name == volumeName || m.Source == volumeName {
				name := c.Names[0]
				if len(name) > 0 && name[0] == '/' {
					name = name[1:]
				}
				usedBy = append(usedBy, name)
				break // Found usage in this container, move to next container
			}
		}
	}
	return usedBy, nil
}

// GetContainersUsingNetwork returns a list of container names that are attached to the specified network ID.
func (cw *ClientWrapper) GetContainersUsingNetwork(networkID string) ([]string, error) {
	containers, err := cw.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var usedBy []string
	for _, c := range containers {
		if c.NetworkSettings != nil {
			for _, net := range c.NetworkSettings.Networks {
				if net.NetworkID == networkID {
					name := c.Names[0]
					if len(name) > 0 && name[0] == '/' {
						name = name[1:]
					}
					usedBy = append(usedBy, name)
					break
				}
			}
		}
	}
	return usedBy, nil
}

// GetContainerStats retrieves the current CPU and memory usage of a container.
func (cw *ClientWrapper) GetContainerStats(id string) (ContainerStats, error) {
	stats, err := cw.client.ContainerStats(context.Background(), id, false)
	if err != nil {
		return ContainerStats{}, err
	}
	defer stats.Body.Close()

	var v types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
		return ContainerStats{}, err
	}

	var cpuPercent float64
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Calculate memory usage
	// MemUsage is v.MemoryStats.Usage - v.MemoryStats.Stats["cache"]
	var memUsage float64
	if v.MemoryStats.Usage > 0 {
		memUsage = float64(v.MemoryStats.Usage)
		if cache, ok := v.MemoryStats.Stats["cache"]; ok {
			memUsage -= float64(cache)
		}
	}

	// Calculate network I/O
	var rx, tx float64
	for _, network := range v.Networks {
		rx += float64(network.RxBytes)
		tx += float64(network.TxBytes)
	}

	return ContainerStats{
		CPUPercent: cpuPercent,
		MemUsage:   memUsage,
		MemLimit:   float64(v.MemoryStats.Limit),
		NetRx:      rx,
		NetTx:      tx,
	}, nil
}

// InspectContainer returns the detailed inspection information for a container.
func (cw *ClientWrapper) InspectContainer(id string) (types.ContainerJSON, error) {
	return cw.client.ContainerInspect(context.Background(), id)
}
