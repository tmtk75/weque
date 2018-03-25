package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/repository"
)

type Handler func(wh repository.Handler, r *http.Request, raw []byte, body *repository.Webhook) error

type Server struct {
	events chan *repository.Webhook
}

func New() *Server {
	return &Server{events: make(chan *repository.Webhook)}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Post("/registry", registry.RegistryHandler)
	r.Post("/repository/github", repository.NewHandler(&repository.Github{}, s.events))
	r.Post("/repository/bitbucket", repository.NewHandler(&repository.Bitbucket{}, s.events))

	go repository.StartConsumer(s.events)

	port := viper.GetInt("port")
	log.Printf("start listening at %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func Validate() error {
	for _, key := range []string{repository.KeyHandlerScript, registry.KeyHandlerScript} {
		path := viper.GetString(key)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
		log.Printf("%v: %v", key, path)
	}
	return nil
}
