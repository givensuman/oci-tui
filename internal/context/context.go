// Package context provides a context for shared application state.
package context

import (
	"sync"

	"github.com/givensuman/containertui/internal/client"
)

var (
	clientInstance *client.ClientWrapper
	once     sync.Once
)

func Init() {
	once.Do(func() {
		clientInstance = client.NewClient()
	})
}

// GetClient returns the shared client instance.
func GetClient() *client.ClientWrapper {
	return clientInstance
}

// CloseClient closes the shared client instance.
func CloseClient() {
	if clientInstance != nil {
		clientInstance.CloseClient()
	}
}
