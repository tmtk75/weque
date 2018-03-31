package registry

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	b, err := ioutil.ReadFile("./payload.json")
	assert.NoError(t, err)

	expects := []struct {
		body   []byte
		expect string
	}{
		{
			body:   b,
			expect: "115d94fb-7318-497c-aa74-4061d4e52548",
		},
		{
			body:   []byte(`{"events":[{"id":"1"},{"id":"2"}]}`),
			expect: "unsupported multiple events: 1,2",
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
	}
}
