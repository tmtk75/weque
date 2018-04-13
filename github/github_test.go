package github

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "abc123")
	req, _ := makeRequest("POST", "/foo/bar", nil)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "api.github.com", req.URL.Host)
	assert.Equal(t, "/foo/bar", req.URL.Path)
	assert.Equal(t, "token abc123", req.Header.Get("authorization"))
}
