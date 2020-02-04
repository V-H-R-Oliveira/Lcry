package main

import (
	"core"
	"os"
)

func main() {
    defer os.RemoveAll(os.Args[0])

    if core.IsRoot() {
		homeDir, _ := os.UserHomeDir()
		core.MountAssets(homeDir)
    }
}
