package repository

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct{}

func (t *Test) Verify(r *http.Request, body []byte) error {
	return nil
}

func (t *Test) IsPing(r *http.Request, body []byte) bool {
	return false
}

func (t *Test) Unmarshal(r *http.Request, body []byte) (*Webhook, error) {
	return &Webhook{}, nil
}

func TestNewHandler(t *testing.T) {
	ch := make(chan<- *Webhook)
	h := NewHandler(&Test{}, ch)
	s := httptest.NewServer(h)
	defer s.Close()

	req := bytes.NewBuffer([]byte("{}"))
	res, err := http.Post(s.URL, "application/json", req)
	assert.NoError(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "", body)
}
