package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

const (
	KeyACMEEnabled       = "tls.acme.enabled"
	KeyACMEPort          = "tls.acme.port"
	KeyACMEHostWhitelist = "tls.acme.host-whitelist"
	KeyACMECacheDir      = "tls.acme.cache-dir"
)

type GetCert func(*tls.ClientHelloInfo) (*tls.Certificate, error)

func ConfigureACME(ch chan<- os.Signal) (*http.Server, GetCert) {
	cachedir := viper.GetString(KeyACMECacheDir)
	whitelist := viper.GetStringSlice(KeyACMEHostWhitelist)

	log.Printf("%v: %v", KeyACMECacheDir, cachedir)
	log.Printf("%v: %v", KeyACMEHostWhitelist, whitelist)

	mgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(whitelist...),
		Cache:      autocert.DirCache(cachedir), // to store certs
	}

	addr := viper.GetString(KeyACMEPort)
	challenge := &http.Server{Addr: addr, Handler: mgr.HTTPHandler(nil)}

	go func() {
		log.Printf("start listening at %v", addr)
		err := challenge.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("http: %v", err)
			ch <- os.Interrupt
		}
	}()

	return challenge, mgr.GetCertificate
}
