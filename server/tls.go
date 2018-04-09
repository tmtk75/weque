package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

const (
	KeyACMEEnabled       = "tls.acme.enabled"
	KeyACMEChallengePort = "tls.acme.challenge-port"
	KeyACMEHostWhitelist = "tls.acme.host-whitelist"
	KeyACMECacheDir      = "tls.acme.cache-dir"
)

func init() {
	viper.SetDefault(KeyACMECacheDir, "certs")
	viper.SetDefault(KeyACMEHostWhitelist, []string{"example.com"})
}

func ListenAndServeTLS(h http.Handler) error {
	cachedir := viper.GetString(KeyACMECacheDir)
	log.Printf("%v: %v", KeyACMECacheDir, cachedir)

	whitelist := viper.GetStringSlice(KeyACMEHostWhitelist)
	log.Printf("%v: %v", KeyACMEHostWhitelist, whitelist)

	mgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(whitelist...),
		Cache:      autocert.DirCache(cachedir), // to store certs
	}

	server := &http.Server{
		Addr:    viper.GetString(KeyTLSPort),
		Handler: h,
		TLSConfig: &tls.Config{
			//GetCertificate: mgr.GetCertificate,
			GetCertificate: GetCertificate,
		},
	}

	addr := viper.GetString(KeyACMEChallengePort)
	challenge := &http.Server{Addr: addr, Handler: mgr.HTTPHandler(nil)}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		log.Printf("start listening at %v", addr)
		err := challenge.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("http: %v", err)
			ch <- os.Interrupt
		}
	}()

	go func() {
		log.Printf("start listening at %v", server.Addr)
		err := server.ListenAndServeTLS("", "") // Coming from GetCertificate
		if err != http.ErrServerClosed {
			log.Printf("https: %v", err)
			ch <- os.Interrupt
		}
	}()

	select {
	case _ = <-ch:
		log.Printf("challenge shutdown: %v", challenge.Shutdown(context.Background()))
		log.Printf("server shutdown: %v", server.Shutdown(context.Background()))
	}
	return nil
}

// for testing
func GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	c, err := tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

//go:generate go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,localhost --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

const certPEMBlock = ` -----BEGIN CERTIFICATE-----
MIICETCCAXqgAwIBAgIQPtOOzqt5KqzFJ+f8uzZ7ijANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDBAIckkd9cq5WrtIVG1uBGMQTr6hLMfDxfA3OwJCblPK1JM97Rl4mek245
y39u7npDvDlmTHJhvF2GK8RxVVLsHtjQ3M5bc7vEQLI6uAD5/wzNFPjNtioeVQnQ
kE6UGjjAuYcXOWdYhNXKoMRZ42BnIQzY9F8EPiWy845ze26SeQIDAQABo2YwZDAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAsBgNVHREEJTAjgglsb2NhbGhvc3SHBH8AAAGHEAAAAAAAAAAAAAAAAAAA
AAEwDQYJKoZIhvcNAQELBQADgYEASe6E+0WXUj54lCudy9jVmFOcs2ILvfsLX2QN
/Ur8VyQbajQQeybva3xA95JQw8yv3rEF4eowCrlQoeOuNBPHCVjMjz0we3oNzgrO
CuLiP5DUMurn8DBVsSceRB0VYvmammYd6/dieLKVGqWtLvEyWKi1ivwWvVFdW1F0
WC3avFw=
-----END CERTIFICATE-----`

const keyPEMBlock = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDBAIckkd9cq5WrtIVG1uBGMQTr6hLMfDxfA3OwJCblPK1JM97R
l4mek245y39u7npDvDlmTHJhvF2GK8RxVVLsHtjQ3M5bc7vEQLI6uAD5/wzNFPjN
tioeVQnQkE6UGjjAuYcXOWdYhNXKoMRZ42BnIQzY9F8EPiWy845ze26SeQIDAQAB
AoGAI0M7bd0RGFdpQzP6XdUIqQpvwcLEqIPSa/Gvg3E3gg6yAnvtrBGp3UVGkFyz
7cq4oAOV4TD6fQzzcX4xqBtUyOnnTIYTCuisbBhrksitMp8a8F0DDjIJjNkiBdWd
h4MlQ7ZednO/qiY93SlTGfuWZEFqPLoxtwSNSl67Vm58XNkCQQDwj6cP29Ia0G98
JwC/6S7HiXQyizqR854ZimnUX9sk+mZyw33FJoA+RP2UAbhvUaOsEEod3HfoJMsq
eykxAnHTAkEAzWN+b0ViDr8EGw2S5zQSnlDcHSYkEnkkzGPt3HxnGivG2aN4OcGN
cugg1Srkhym/fAgDsJ2htbwr5ory2givAwJBANKKB377tuE897XDNQbBgO2mQYpT
DIncm8xitcjntBajCLL8ocDAt5DINN8quk7DNupKv3NvF4qXWTDu5dg8+X0CQQCR
zZun6h1eUoPboJs0vmapNMXNe5IH+zAAWMA20alvjrwvLDjg52I+vELykOyCd0SU
DCxyaLSvitGva9xSo+95AkEAt5T8hQykomlgWSXntbQu/CEvMqKfB3o+Tw35VmYu
8pNBxfmkUsWTA1ma4AHJ9g7z0ctL7UG2ycuEJCRtYy82XQ==
-----END RSA PRIVATE KEY-----`
