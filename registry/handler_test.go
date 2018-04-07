package registry

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	b, err := ioutil.ReadFile("./payload.json")
	assert.NoError(t, err)

	expects := []struct {
		body   []byte
		expect string
		code   int
	}{
		{
			body:   b,
			expect: "115d94fb-7318-497c-aa74-4061d4e52548",
			code:   200,
		},
		{
			body:   []byte(`{"events":[{"id":"1"},{"id":"2"}]}`),
			expect: "unsupported multiple events: 1,2",
			code:   400,
		},
		{
			body:   []byte(`{}`),
			expect: "no any events",
			code:   400,
		},
	}

	s := httptest.NewServer(NewHandler(make(chan<- *Webhook)))
	defer s.Close()

	for _, e := range expects {
		res, err := http.Post(s.URL, "application/json", bytes.NewBuffer(e.body))
		assert.NoError(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, e.expect, string(body))
		assert.Equal(t, e.code, res.StatusCode)
	}
}

func TestHandlerVeirfy(t *testing.T) {
	expects := []struct {
		reqbody string
		headers map[string]string
		code    int
		body    string
	}{
		{
			reqbody: "{}",
			headers: map[string]string{},
			code:    400,
			body:    "failed to verify: the given secret token didn't match",
		},
		{
			reqbody: `{"events":[{"id":"foobar"}]}`,
			headers: map[string]string{"x-weque-secret": "abc123"},
			code:    200,
			body:    "foobar",
		},
	}

	viper.Set(KeySecretToken, "abc123")

	t.Run("test", func(t *testing.T) {
		for _, e := range expects {
			req := httptest.NewRequest("POST", "http://example.com/registry", bytes.NewBuffer([]byte(e.reqbody)))
			for k, v := range e.headers {
				req.Header.Add(k, v)
			}

			w := httptest.NewRecorder()
			h := NewHandler(make(chan<- *Webhook))
			h(w, req)

			res := w.Result()
			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, e.code, res.StatusCode)
			assert.Equal(t, e.body, string(body))
		}
	})

	viper.Set(KeySecretToken, "")
}

func TestHandlerInsecure(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com", nil)
	h := NewHandler(make(chan<- *Webhook))
	viper.Set(KeySecretToken, "abc123")

	// secure
	viper.Set(KeyInsecureMode, false)
	t.Run("secure", func(t *testing.T) {
		r := httptest.NewRecorder()
		h(r, req)
		res := r.Result()
		assert.Equal(t, 400, res.StatusCode)
	})

	// insecure
	viper.Set(KeyInsecureMode, true)
	t.Run("insecure", func(t *testing.T) {
		r := httptest.NewRecorder()
		h(r, req)
		res := r.Result()
		assert.Equal(t, 200, res.StatusCode)
	})

	viper.Set(KeyInsecureMode, false)
}
