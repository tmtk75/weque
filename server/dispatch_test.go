package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
)

func init() {
	viper.Set(repository.KeySecretToken, "abc123")
}

func TestNewDispatcher(t *testing.T) {
	var (
		ch        = make(chan<- *repository.Context)
		github    = repository.NewHandler(&repository.Github{}, ch)
		bitbucket = repository.NewHandler(&repository.Bitbucket{}, ch)
	)
	h := NewDispatcher(github, bitbucket)

	s := httptest.NewServer(h)
	defer s.Close()

	data := []struct {
		src, query string
		headers    map[string]string
		expected   string
	}{
		{
			src: "../github/payload.json", query: "",
			headers: map[string]string{
				"Content-Type":      "application/json",
				"X-Hub-Signature":   "sha1=c699905923f6a533824e8fb13a0b344d52146e20",
				"X-Github-Delivery": "gh-json",
			},
			expected: "gh-json\n",
		},
		{
			src: "../github/payload.txt", query: "",
			headers: map[string]string{
				"Content-Type":      "application/x-www-form-urlencoded",
				"X-Hub-Signature":   "sha1=af9c4634ebadf38f19f14c713f2ab9c0328934ad",
				"X-Github-Delivery": "gh-urlencoded",
			},
			expected: "gh-urlencoded\n",
		},
		{
			src: "../bitbucket/payload.json", query: "?secret=abc123", // secret is provided by query parameter for bitbucket
			headers: map[string]string{
				"Content-Type":   "application/json",
				"X-Request-UUID": "bb-foobar",
			},
			expected: "bb-foobar\n",
		},
	}

	for _, e := range data {
		p, err := ioutil.ReadFile(e.src)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", s.URL+e.query, bytes.NewBuffer(p))
		for k, v := range e.headers {
			req.Header.Add(k, v)
		}
		res, err := (&http.Client{}).Do(req)

		assert.NoError(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, e.expected, string(body))
	}
}
