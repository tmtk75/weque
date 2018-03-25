package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/consumer"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/repository"
)

type Server struct {
	events chan *repository.Context
}

func New() *Server {
	return &Server{events: make(chan *repository.Context)}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Post("/registry", registry.RegistryHandler)
	r.Post("/repository/github", repository.NewHandler(&repository.Github{}, s.events))
	r.Post("/repository/bitbucket", repository.NewHandler(&repository.Bitbucket{}, s.events))

	go consumer.StartRepositoryConsumer(s.events)

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
	for _, key := range []string{consumer.KeyHandlerScriptRepository, registry.KeyHandlerScript} {
		path := viper.GetString(key)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
		log.Printf("%v: %v", key, path)
	}
	return nil
}