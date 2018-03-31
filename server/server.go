package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/consumer"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/repository"
)

type Server struct {
	repositoryEvents chan *repository.Context
	registryEvents   chan *registry.Webhook
}

func New() *Server {
	return &Server{
		repositoryEvents: make(chan *repository.Context),
		registryEvents:   make(chan *registry.Webhook),
	}
}

func (s *Server) Start() error {
	var (
		gh        = &repository.Github{}
		bb        = &repository.Bitbucket{}
		github    = repository.NewHandler(gh, s.repositoryEvents)
		bitbucket = repository.NewHandler(bb, s.repositoryEvents)
		regh      = NewDispatcher(github, bitbucket)
	)
	e := echo.New()
	e.POST("/registry", Wrap(e, registry.NewHandler(s.registryEvents)))
	e.POST("/", Wrap(e, regh))
	e.POST("/repository", Wrap(e, regh))
	e.POST("/repository/github", Wrap(e, github))
	e.POST("/repository/bitbucket", Wrap(e, bitbucket))

	go consumer.StartRepositoryConsumer(s.repositoryEvents)

	port := viper.GetInt("port")
	log.Printf("start listening at %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), e)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func Wrap(e *echo.Echo, h http.HandlerFunc) func(echo.Context) error {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
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
