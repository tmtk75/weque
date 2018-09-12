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
	gl "github.com/tmtk75/weque/repository/gitlab"
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
		gitlab    = repository.NewHandler(&gl.Gitlab{}, s.repositoryEvents)
		regh      = NewDispatcher(github, bitbucket)
	)
	e := echo.New()
	e.POST("/registry", Wrap(e, registry.NewHandler(s.registryEvents)))
	e.POST("/", Wrap(e, regh))
	e.POST("/repository", Wrap(e, regh))
	e.POST("/repository/github", Wrap(e, github))
	e.POST("/repository/bitbucket", Wrap(e, bitbucket))
	e.POST("/repository/gitlab", Wrap(e, gitlab))
	// Deprecated. To be compatible with https://github.com/tmtk75/hoko
	e.POST("/serf/event/github", Wrap(e, github))
	e.POST("/serf/event/bitbucket", Wrap(e, bitbucket))

	go printError(repositoryworker.Notify(repositoryworker.Run(s.repositoryEvents)))
	go printError(registryworker.Notify(registryworker.Run(s.registryEvents)))

	//log.Printf("insecure: %v", viper.GetBool(KeyInsecure))

	var err error
	if viper.GetBool(KeyTLSEnabled) {
		err = ListenAndServeTLS(e)
	} else {
		err = ListenAndServe(e)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func printError(ch <-chan error) {
	for err := range ch {
		if err != nil {
			log.Printf("[error] %v", err)
		}
	}
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
