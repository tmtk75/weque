package weque

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Request(mkreq func(m, p string, b io.Reader) (*http.Request, error), method, path string, body io.Reader) (string, error) {
	req, err := mkreq(method, path, body)
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
