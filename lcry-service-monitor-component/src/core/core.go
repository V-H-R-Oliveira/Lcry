package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"networking"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// getServicePid -> get lcry.service PID
func getServicePid(service string) (string, error) {
	cmd := exec.Command("/bin/pidof", service, "echo stdout; echo 1>&2 stderr")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s was removed from the system...\n", service)
		return "", err
	}

	lcryPid := fmt.Sprintf(strings.Trim(string(output), "\n\r\t"))
	return lcryPid, nil
}

// getServiceStatus -> gets the current status of the service
func getServiceStatus(pid string) string {
	procInfoPath := filepath.Join("/proc", pid, "stat")

	procInfo, err := ioutil.ReadFile(procInfoPath)

	if err != nil {
		log.Fatal(err)
	}

	infoVector := strings.Split(string(procInfo), " ")
	return infoVector[2]
}

// StartMonitor -> Starts the service monitor
func StartMonitor(service, targetName string) {
	for {
		pid, err := getServicePid(service)

		if err != nil {
			networking.SendRes(fmt.Sprintf("[From %s] monitoring %s, and the status is %s\n", targetName, service, "D"))
			return
		}

		switch status := getServiceStatus(pid); status {
		case "R":
			networking.SendRes(fmt.Sprintf("[From %s] monitoring %s, and the status is %s\n", targetName, service, "R"))
			break
		case "S":
			networking.SendRes(fmt.Sprintf("[From %s] monitoring %s, and the status is %s\n", targetName, service, "S"))
			break
		default:
			networking.SendRes(fmt.Sprintf("[From %s] monitoring %s, and the status is %s\n", targetName, service, "U"))
			return
		}

		time.Sleep(3 * time.Minute)
	}
}
