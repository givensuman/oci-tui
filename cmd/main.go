package main

import (
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui"
)

func main() {
	context.Init()
	defer context.CloseClient()

	ui.Start()
}
