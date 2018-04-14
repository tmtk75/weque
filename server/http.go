package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

func ListenAndServe(e http.Handler) error {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	port := viper.GetString(KeyPort)
	log.Printf("start listening at %s", port)

	server := &http.Server{Addr: port, Handler: e}
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("http: %v", err)
			ch <- os.Interrupt
		}
	}()

	_ = <-ch
	log.Printf("http shutdown: %v", server.Shutdown(context.Background()))

	return nil
}
