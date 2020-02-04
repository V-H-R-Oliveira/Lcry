package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

func verifyMagic(filePath string) bool {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	magic := make([]byte, 9)

	_, err = file.Read(magic)

	if err != nil {
		log.Fatal(err)
	}

	checksumtoHex := fmt.Sprintf("%x", magic)

	if cmp := strings.Compare(checksumtoHex, "6c6372797074353072"); cmp != 0 {
		return false
	}

	return true
}

func decryptDir(dirPath string, key *[]byte) {
	ch := make(chan string)

	go func() {
		e := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if !info.IsDir() && verifyMagic(path) {
				ch <- path
			}

			return nil
		})

		close(ch)

		if e != nil {
			log.Fatal(e)
		}
	}()

	decryptWrapper(ch, key)
}

func decryptWrapper(ch chan string, key *[]byte) {
	for path := range ch {
		decrypt(path, key)
	}
}

const bufferSize = 10485760

func decrypt(filePath string, key *[]byte) {
	if len(*key) != 32 {
		log.Fatal("Nop")
	}

	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	defer file.Close()

	endPath := strings.Trim(filePath, filepath.Ext(filePath))
	decFile, err := os.Create(endPath)

	if err != nil {
		log.Fatal(err)
	}

	defer decFile.Close()

	magic := make([]byte, 9)

	_, err = file.Read(magic)

	if err != nil {
		log.Fatal(err)
	}

	encodeToHex := fmt.Sprintf("%x", magic)

	if cmp := strings.Compare(encodeToHex, "6c6372797074353072"); cmp != 0 {
		log.Fatal("Magic not present.")
	}

	iv := make([]byte, aes.BlockSize)

	_, err = file.Read(iv)

	if err != nil {
		log.Fatal(err)
	}

	hmac := hmac.New(sha256.New, *key)
	hmac.Write(iv)

	fileinfo, err := os.Stat(filePath)

	if err != nil {
		log.Fatal(err)
	}

	offset := fileinfo.Size() - 32
	checksum := make([]byte, 32)
	_, err = file.ReadAt(checksum, offset)

	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(*key)

	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCTR(block, iv)
	buffer := make([]byte, bufferSize)

	for {
		nBytes, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			break
		}

		plaintext := make([]byte, nBytes)

		if nBytes == bufferSize {
			hmac.Write(buffer[:nBytes])
			stream.XORKeyStream(plaintext, buffer[:nBytes])
			decFile.Write(plaintext)
		} else if nBytes < bufferSize && nBytes != 32 {
			limit := nBytes - 32
			hmac.Write(buffer[:limit])
			stream.XORKeyStream(plaintext, buffer[:limit])
			plaintext = plaintext[:limit]
			decFile.Write(plaintext)
		} else {
			break
		}
	}

	hash := hmac.Sum(nil)
	authEnc(endPath, &checksum, &hash)
	os.Remove(filePath)
}

func authEnc(endPath string, checksum, extractedChecksum *[]byte) {
	if len(*extractedChecksum) != len(*checksum) {
		fmt.Println("Different Checksum length...")
		log.Fatal(os.Remove(endPath))
	}

	for index, bytee := range *checksum {
		if bytee != (*extractedChecksum)[index] {
			fmt.Println("Bytes mismatch - different checkum code...")
			log.Fatal(os.Remove(endPath))
		}
	}
}

func extractPbkdf2Key(password, salt *[]byte) []byte {
	const iter = 64000
	const keySize = 32

	return pbkdf2.Key(*password, *salt, iter, keySize, sha3.New384)
}

func main() {
	var (
		passwd, salt, path string
	)

	flag.StringVar(&passwd, "passwd", "", "Original password for decrypting the files.")
	flag.StringVar(&salt, "salt", "", "Original salt in hexadecimal format.")
	flag.StringVar(&path, "path", "", "Decrypt path")
	flag.Parse()

	if flag.Parsed() {
		if passwd != "" && salt != "" && path != "" {
			salt, _ := hex.DecodeString(salt)
			passwdBytes := []byte(passwd)
			usedKey := extractPbkdf2Key(&passwdBytes, &salt)
			decryptDir(path, &usedKey)
		} else {
			log.Fatal("Flags usage: -passwd=<your passwd> -salt=<hex salt> -path=<path to be decrypted>")
		}
	}
}
