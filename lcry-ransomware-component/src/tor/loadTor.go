package tor

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

const (
	mfdCloexec  = 0x1 // flag
	memfdCreate = 319 // syscall number
)

// RunTor -> Run embbed Tor
func RunTor(target string) {
	fdName := ""
	fd, _, _ := syscall.Syscall(memfdCreate, uintptr(unsafe.Pointer(&fdName)), uintptr(mfdCloexec), 0)

	_, err := syscall.Write(int(fd), *embbedTor)

	if err != nil {
		log.Fatal(err)
	}

	fdPath := fmt.Sprintf("/proc/self/fd/%d", fd)
	argvs := []string{fdPath}
	env := []string{"HOME=" + target}

	if _, err = syscall.ForkExec(fdPath, argvs, &syscall.ProcAttr{Env: env}); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
    return
}
