package bitbucket

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitbucketWebhookProvider(t *testing.T) {
	p, err := ioutil.ReadFile("../../bitbucket/payload.json")
	assert.NoError(t, err)

	b := &Bitbucket{}
	r := httptest.NewRequest("POST", "/", nil)

	w, err := b.Unmarshal(r, p)
	assert.NoError(t, err)

	L := "https://bitbucket.org/tmtk"
	R := "knockout.meteor-demo"
	assert.Equal(t, L+"/"+R+"", b.RepositoryURL(w))
	assert.Equal(t, L+"/"+R+"/commits/2b24041b4efd416fa30109da37c6f025696aa92c", b.CommitURL(w))
	assert.Equal(t, L+"/"+R+"/branches/compare/c33b60b820b104d176df6eb6f79a187f972c943f..2b24041b4efd416fa30109da37c6f025696aa92c", b.CompareURL(w))
	assert.Equal(t, L+"/"+R+"/branch/master", b.RefURL(w))
	assert.Equal(t, L, b.PusherURL(w))
}
