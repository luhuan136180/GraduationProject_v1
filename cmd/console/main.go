package main

import (
	"go.uber.org/zap"
	console "v1/cmd/console/app"
)

func main() {
	cmd := console.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		zap.S().Fatal(err)
	}
}
