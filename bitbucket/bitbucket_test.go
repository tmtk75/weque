package bitbucket

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	os.Setenv("BITBUCKET_API_KEY", "abc123")
	req, _ := makeRequest("user", "POST", "/foo/bar", nil)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "api.bitbucket.org", req.URL.Host)
	assert.Equal(t, "/2.0/repositories/foo/bar", req.URL.Path)
	assert.Equal(t, "Basic dXNlcjphYmMxMjM=", req.Header.Get("authorization"))
}
