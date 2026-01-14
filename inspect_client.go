package main

import (
	"fmt"
	"reflect"

	"github.com/moby/moby/client"
)

func main() {
	c, _ := client.New(client.FromEnv)
	t := reflect.TypeOf(c)
	fmt.Printf("Client Type: %s\n", t)
	for i := range t.NumMethod() {
		m := t.Method(i)
		fmt.Printf("Method: %s\n", m.Name)
	}
	fmt.Println("Scan complete")
}
