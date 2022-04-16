package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func ProjectRoot() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		var err error
		thisFile, err = os.Getwd()
		if err != nil {
			panic(errors.Wrapf(err, "cant determine the root path"))
		}
	}
	projectRoot, _ := filepath.Abs(fmt.Sprintf("%s/../../", filepath.Dir(thisFile)))
	return projectRoot
}

func LoadEnv() error {
	envPath := filepath.Join(ProjectRoot(), ".env")
	return godotenv.Load(envPath)
}
