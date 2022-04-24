package main

import (
	"github.com/reliefeffortslk/protest-tracker-api/cmd"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/configs"
)

var _ = configs.LoadEnv()

func main() {
	cmd.Execute()
}
