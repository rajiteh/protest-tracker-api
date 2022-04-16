package main

import (
	"github.com/reliefeffortslk/protest-tracker-api/cmd"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/configs"
)

var _ = configs.LoadEnv()

func main() {
	cmd.Execute()
}

// https://docs.google.com/spreadsheets/d/1yShvemHd_eNNAtC3pmxPs9B5RbGmfBUP1O6WGQ5Ycrg/pub?output=csv
// https://docs.google.com/spreadsheets/d/e/2PACX-1vQIhLNNfUKVjxMkMwdtTFnvuV8oN1H_OmgOWRCwHBkSfOo1fzA08LXDfcK4EA86fx18M4FeAIwOoBBR/pub?output=csv
