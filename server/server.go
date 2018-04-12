package server

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/registry/worker"
	"github.com/tmtk75/weque/repository"
	bb "github.com/tmtk75/weque/repository/bitbucket"
	gh "github.com/tmtk75/weque/repository/github"
	"github.com/tmtk75/weque/repository/worker"
)

const (
	KeyPort    = "port"
	KeyTLSPort = "tls.port"
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
		github    = repository.NewHandler(&gh.Github{}, s.repositoryEvents)
		bitbucket = repository.NewHandler(&bb.Bitbucket{}, s.repositoryEvents)
		regh      = NewDispatcher(github, bitbucket)
	)
	e := echo.New()
	e.POST("/registry", Wrap(e, registry.NewHandler(s.registryEvents)))
	e.POST("/", Wrap(e, regh))
	e.POST("/repository", Wrap(e, regh))
	e.POST("/repository/github", Wrap(e, github))
	e.POST("/repository/bitbucket", Wrap(e, bitbucket))

	repositoryworker.Notify(repositoryworker.Run(s.repositoryEvents))
	registryworker.Notify(registryworker.Run(s.registryEvents))

	var err error
	if !viper.GetBool(KeyACMEEnabled) {
		err = ListenAndServe(e)
	} else {
		err = ListenAndServeTLS(e)
	}
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
	for _, key := range []string{repositoryworker.KeyHandlerScriptRepository, registryworker.KeyHandlerScriptRegistry} {
		path := viper.GetString(key)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
		log.Printf("%v: %v", key, path)
	}
	return nil
}

func ListenAndServe(e http.Handler) error {
	port := viper.GetString(KeyPort)
	log.Printf("start listening at %s", port)
	err := http.ListenAndServe(port, e)
	return err
}
