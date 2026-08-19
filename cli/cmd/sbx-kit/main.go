package main

import (
	"os"

	"github.com/nkapatos/sbx-kit/cli/internal/cmd"
)

func main() {
	if err := cmd.NewRoot().Execute(); err != nil {
		cmd.UI().ErrorPrefix(err.Error())
		os.Exit(1)
	}
}
