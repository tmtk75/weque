package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

const (
	KeyTLSEnabled  = "tls.enabled"
	KeyTLSCertFile = "tls.cert_file"
	KeyTLSKeyFile  = "tls.key_file"
)

func init() {
	viper.SetDefault(KeyACMECacheDir, "certs")
	viper.SetDefault(KeyACMEHostWhitelist, []string{"test.example.com"})
}

type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

type NopShutdownable struct{}

func (e *NopShutdownable) Shutdown(ctx context.Context) error {
	log.Println("[debug] nop shutdownable")
	return nil
}

func ListenAndServeTLS(h http.Handler) error {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	var challenge Shutdownable = &NopShutdownable{}
	getCert := GetCertificateLocalhost
	if viper.GetBool(KeyACMEEnabled) {
		challenge, getCert = ConfigureACME(ch)
	} else {
		getCert = GetCertificateFile
	}

	server := &http.Server{
		Addr:    viper.GetString(KeyTLSPort),
		Handler: h,
		TLSConfig: &tls.Config{
			GetCertificate: getCert,
		},
	}

	go func() {
		log.Printf("start listening at %v", server.Addr)
		err := server.ListenAndServeTLS("", "") // Coming from GetCertificate
		if err != http.ErrServerClosed {
			log.Printf("https: %v", err)
			ch <- os.Interrupt
		}
	}()

	_ = <-ch
	log.Printf("challenge shutdown: %v", challenge.Shutdown(context.Background()))
	log.Printf("tls shutdown: %v", server.Shutdown(context.Background()))
	return nil
}

func GetCertificateFile(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	var (
		cert = viper.GetString(KeyTLSCertFile)
		key  = viper.GetString(KeyTLSKeyFile)
	)
	log.Printf("cert_file: %v, key_file: %v", cert, key)
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// for testing
func GetCertificateLocalhost(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	c, err := tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

//#go:generate go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,localhost --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

const certPEMBlock = `-----BEGIN CERTIFICATE-----
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
