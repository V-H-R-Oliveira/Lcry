package core

import (
    "fmt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
	"os"
)

const bufferSize = 10485760

func encrypt(filePath string, key *[]byte) {
	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	defer file.Close()
	defer os.Remove(filePath)

	endPath := filePath + ".lcry"

	encFile, err := os.Create(endPath)

	if err != nil {
		return
	}

	defer encFile.Close()

	encFile.WriteString("lcrypt50r")

	block, err := aes.NewCipher(*key)

	if err != nil {
		fmt.Println("Core err 1")
        log.Fatal(err)
	}

	buffer := make([]byte, bufferSize)
	iv := make([]byte, aes.BlockSize)
	hmac := hmac.New(sha256.New, *key)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	encFile.Write(iv)
	hmac.Write(iv)

	stream := cipher.NewCTR(block, iv)

	for {
		nBytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Core err 2")
            log.Fatal(err)
		}

		if err == io.EOF {
			break
		}

		ciphertext := make([]byte, nBytes)
		stream.XORKeyStream(ciphertext, buffer[:nBytes])
		hmac.Write(ciphertext)
		encFile.Write(ciphertext)
	}

	hash := hmac.Sum(nil)
	encFile.Write(hash)
}
