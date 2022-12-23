package workspace

import (
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"gopkg.in/yaml.v3"
)

var (
	pathSeparator = string(os.PathSeparator)
)

func CheckWorkSpace(dir string) []string {
	main := filepath.Join(dir, "package.json")
	if !exist(main) {
		return []string{}
	}
	if isWorkSpace(dir) {
		array := []string{main}
		config := getWorkspaceConfig(dir)
		for _, subdir := range config.Packages {
			path := filepath.Join(dir, subdir)
			glob, _ := zglob.Glob(path + pathSeparator + "**" + pathSeparator + "package.json")
			array = append(array, glob...)
		}
		return array
	}
	return []string{main}
}

func isWorkSpace(dir string) bool {
	file := filepath.Join(dir, "pnpm-workspace.yaml")
	return exist(file)
}

func exist(path string) bool {
	glob, err := filepath.Glob(path)
	if err != nil {
		return false
	}
	return len(glob) > 0
}

func getWorkspaceConfig(dir string) Config {
	var config Config
	file := filepath.Join(dir, "pnpm-workspace.yaml")
	byte, _ := os.ReadFile(file)
	yaml.Unmarshal(byte, &config)
	return config
}
