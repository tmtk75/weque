package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/repository"
	bb "github.com/tmtk75/weque/repository/bitbucket"
	gh "github.com/tmtk75/weque/repository/github"
)

func NewDispatcher(github, bitbucket http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			weque.SendError(w, 500, fmt.Sprintf("%v", err))
			return
		}
		log.Printf("dispatcher read %vbytes", len(b))
		//log.Printf("[debug] %v...", string(b)[0:10])

		var wh *repository.Webhook
		// GitHub
		wh, err = (&gh.Github{}).Unmarshal(r, b)
		if err != nil {
			// Not github
			log.Printf("[info] failed to parse request body as github: %v", err)
		} else if wh.After != "" {
			log.Printf("[info] was able to parse as github")
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			github(w, r)
			return
		} else if r.Header.Get("x-github-event") == "ping" {
			log.Printf("[info] was able to parse as github ping")
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			github(w, r)
			return
		} else {
			log.Printf("missing After. It seems not valid body thought parsed. %v", wh)
		}
		//log.Printf("after: %v", wh.After)

		// Bitbucket
		wh, err = (&bb.Bitbucket{}).Unmarshal(r, b)
		if err != nil {
			// Not bitbucket
			log.Printf("[info] failed to parse request body as bitbucket: %v", err)
		} else if wh.After != "" {
			log.Printf("[info] was able to parse as bitbucket")
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			bitbucket(w, r)
			return
		} else {
			log.Printf("missing After. It seems not valid body thought parsed. %v", wh)
		}

		// Unknown
		weque.SendError(w, 400, "unknown webhook: failed to unmarshal as github and bitbucket.")
	}
}
