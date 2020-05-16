package main

import (
	"github.com/tiandi111/raft/cmd/app"
	"os"
)

func main() {
	if err := app.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
