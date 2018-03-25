package repository

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubWebhookProvider(t *testing.T) {
	p, err := ioutil.ReadFile("../github/payload.json")
	assert.NoError(t, err)

	b := &Github{}
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Add("content-type", "application/json")

	w, err := b.Unmarshal(r, p)
	assert.NoError(t, err)

	L := "https://github.com/tmtk75"
	R := "weque"
	assert.Equal(t, L+"/"+R+"", b.RepositoryURL(w))
	assert.Equal(t, L+"/"+R+"/commit/bef7731a26b4546fd87dfa73b8a94d9b4f22985e", b.CommitURL(w))
	assert.Equal(t, L+"/"+R+"/compare/f4aa71ad25dd499eeca1a2c6223a7dd232af27bf...bef7731a26b4546fd87dfa73b8a94d9b4f22985e?w=1", b.CompareURL(w))
	assert.Equal(t, L+"/"+R+"/tree/refs/heads/master", b.RefURL(w))
	assert.Equal(t, L, b.PusherURL(w))
}
