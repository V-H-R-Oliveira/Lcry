package main

import (
	"path/filepath"
	"os"
	"core"
)

const service = "winner.lcry"

func main() {
	targetPath, _ := os.UserHomeDir()
	targetName := filepath.Base(targetPath)
	core.StartMonitor(service, targetName)
}
