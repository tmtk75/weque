package repository

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type Test struct {
}

func (t *Test) Verify(r *http.Request, body []byte) error {
	if "expect to fail" == string(body) {
		return errors.New("hi")
	}
	return nil
}

func (t *Test) IsPing(r *http.Request, body []byte) bool {
	return false
}

func (t *Test) Unmarshal(r *http.Request, body []byte) (*Webhook, error) {
	return &Webhook{}, nil
}

func (t *Test) RequestID(r *http.Request) string {
	return ""
}

func (t *Test) WebhookProvider() WebhookProvider {
	return nil
}

func TestNewHandler(t *testing.T) {
	ch := make(chan<- *Context)
	h := NewHandler(&Test{}, ch)
	s := httptest.NewServer(h)
	defer s.Close()

	req := bytes.NewBuffer([]byte("{}"))
	res, err := http.Post(s.URL, "application/json", req)
	assert.NoError(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	expected := []uint8([]byte{0xa})
	assert.Equal(t, expected, body)
}

func TestHandlerInsecure(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com", bytes.NewBufferString("expect to fail"))
	h := NewHandler(&Test{}, make(chan<- *Context))
	viper.Set(KeySecretToken, "abc123")

	// secure
	viper.Set(KeyInsecureMode, false)
	t.Run("secure", func(t *testing.T) {
		r := httptest.NewRecorder()
		h(r, req)
		res := r.Result()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, 401, res.StatusCode)
		assert.Regexp(t, "failed to verify", string(body))
	})

	// insecure
	viper.Set(KeyInsecureMode, true)
	t.Run("insecure", func(t *testing.T) {
		r := httptest.NewRecorder()
		h(r, req)
		res := r.Result()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, 200, res.StatusCode)
		assert.Regexp(t, "", string(body))
	})

	viper.Set(KeyInsecureMode, false)
}
