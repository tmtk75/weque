package gitlab

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tmtk75/weque"
)

/*
 * repo: ex) tmtk75/foobar
 */
func List(repo string) {
	path := fmt.Sprintf("/projects/%s/hooks", url.PathEscape(repo))
	//log.Println(path)
	s, err := Request("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func Create(repo, uri, secret string) {
	path := fmt.Sprintf("/projects/%s/hooks", url.PathEscape(repo))

	d := url.Values{}
	d.Set("url", uri)
	d.Add("token", secret)

	b := bytes.NewBufferString(d.Encode())
	s, err := Request("POST", path, b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func Delete(repo, id string) {
	path := fmt.Sprintf("/projects/%s/hooks/%v", url.PathEscape(repo), id)

	s, err := Request("DELETE", path, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func Request(method, path string, body io.Reader) (string, error) {
	return weque.Request(makeRequest, method, path, body)
}

func makeRequest(method, path string, body io.Reader) (*http.Request, error) {
	var (
		token    = os.Getenv("GITLAB_PRIVATE_TOKEN")
		endpoint = "https://gitlab.com/api/v4"
	)

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", endpoint, path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	return req, nil
}
