package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

func getExecAbsolutePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index]
}

var (
	execPath    = getExecAbsolutePath()
	__DEV__     = strings.Contains(execPath, "Temp") || strings.Contains(execPath, "tmp") || strings.Contains(execPath, "var/folders")
	currentPath string
)

func init() {
	pterm.DefaultSection.Printfln("Fresh!")
	cwd, _ := os.Getwd()
	if __DEV__ {
		result, _ := pterm.DefaultInteractiveSelect.
			WithOptions([]string{"repo", "monorepo"}).
			Show("Which one to debug?")
		currentPath = filepath.Join(cwd, "test", result)
	} else {
		currentPath = cwd
	}
}
