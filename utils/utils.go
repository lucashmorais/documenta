package utils

import (
	"os"
	"path/filepath"
)

func GetServerPath() string {
	env_path, env_path_is_set := os.LookupEnv("DOCUMENTA_ROOT")
	if env_path_is_set {
		return env_path
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		return exPath
	}
}