package bitbucket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func List(repo string) {
	user := strings.Split(repo, "/")[0]
	s, err := Request(user, "GET", fmt.Sprintf("/%s/hooks", repo), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func Create(repo, url, secret string) {
	a := Webhook{
		URL:    url,
		Active: true,
		Events: []string{"repo:push"},
	}
	user := strings.Split(repo, "/")[0]
	s, err := Request(user, "POST", fmt.Sprintf("/%v/hooks", repo), bytes.NewBuffer(a.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

type Webhook struct {
	URL    string   `json:"url"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
}

func (w Webhook) Bytes() []byte {
	b, _ := json.Marshal(w)
	return b
}

func Request(user, method, path string, body io.Reader) (string, error) {
	req, err := makeRequest(user, method, path, body)
	if err != nil {
		return "", err
	}

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return "", err
	}
	if res.StatusCode/100 != 2 {
		log.Print(res.Status)
		return "", errors.Errorf("got %v", res.Status)
	}

	return string(b), nil
}

func makeRequest(user, method, path string, body io.Reader) (*http.Request, error) {
	var (
		apikey   = os.Getenv("BITBUCKET_API_KEY")
		endpoint = "https://api.bitbucket.org/2.0/repositories"
	)

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", endpoint, path), body)
	if err != nil {
		return nil, err
	}
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, apikey)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", token))

	return req, nil
}
