package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.client == nil {
		t.Error("client.client is nil")
	}
}

func TestContainer(t *testing.T) {
	// Test Container struct
	c := Container{
		ID:    "test-id",
		Name:  "test-name",
		Image: "test-image",
		State: "running",
	}
	if c.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", c.ID)
	}
	if c.Name != "test-name" {
		t.Errorf("expected Name test-name, got %s", c.Name)
	}
	if c.Image != "test-image" {
		t.Errorf("expected Image test-image, got %s", c.Image)
	}
	if c.State != "running" {
		t.Errorf("expected State running, got %s", c.State)
	}
}

func TestGetContainers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	client := NewClient()
	defer client.CloseClient()

	containers := client.GetContainers()
	_ = containers
}

func TestGetContainerState(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	client := NewClient()
	defer client.CloseClient()

	// Test with nonexistent container
	state := client.GetContainerState("nonexistent")
	if state != "unknown" {
		t.Errorf("expected 'unknown' for nonexistent container, got %s", state)
	}
}
