// Package client defines means of communicating
// with the OCI runtime
package client

import (
	"log"

	"github.com/docker/docker/client"
)

type Client struct {
	Client *client.Client
}

func NewClient() (*Client, error) {
	var err error
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("failed to create client")
		return nil, err
	}

	return &Client{Client: cli}, nil
}
