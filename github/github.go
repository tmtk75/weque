package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tmtk75/weque"
)

func List(repo string) {
	s, err := Request("GET", fmt.Sprintf("/repos/%s/hooks", repo), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func Create(repo, url, secret, name string) {
	a := Webhook{
		Name:   name,
		Active: true,
		Events: []string{"push"},
		Config: Config{
			URL:         url,
			ContentType: "application/json",
			Secret:      secret,
		},
	}
	s, err := Request("POST", fmt.Sprintf("/repos/%v/hooks", repo), bytes.NewBuffer(a.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

type Webhook struct {
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
	Config Config   `json:"config"`
}

type Config struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Secret      string `json:"secret"`
}

func (w Webhook) Bytes() []byte {
	b, _ := json.Marshal(w)
	return b
}

func Request(method, path string, body io.Reader) (string, error) {
	return weque.Request(makeRequest, method, path, body)
}

func makeRequest(method, path string, body io.Reader) (*http.Request, error) {
	var (
		token    = os.Getenv("GITHUB_TOKEN")
		endpoint = "https://api.github.com"
	)

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", endpoint, path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %v", token))

	return req, nil
}
