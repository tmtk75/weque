package gitlab

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

func Request(method, path string, body io.Reader) (string, error) {
	req, err := makeRequest(method, path, body)
	if err != nil {
		log.Print(err)
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

	return string(b), nil
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
