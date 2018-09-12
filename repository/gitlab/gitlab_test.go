package gitlab

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlab(t *testing.T) {
	p, err := ioutil.ReadFile("../../gitlab/payload.json")
	assert.NoError(t, err)

	b := &Gitlab{}
	r := httptest.NewRequest("POST", "/", nil)

	w, err := b.Unmarshal(r, p)
	assert.NoError(t, err)

	t.Run("Unmarshal", func(t *testing.T) {
		assert.Equal(t, "8c53a97bc3491c125baaacb4974aab5a0ca62c8b", w.After)
		assert.Equal(t, "0000000000000000000000000000000000000000", w.Before)
		assert.Equal(t, false, w.Created) // always false
		assert.Equal(t, false, w.Deleted) // always false
		assert.Equal(t, "tmtk75", w.Pusher.Name)
		assert.Equal(t, "", w.Event)    // always empty
		assert.Equal(t, "", w.Delivery) // always empty
		assert.Equal(t, "refs/heads/0.0.8", w.Ref)
		assert.Equal(t, "google-sheets-sample-go", w.Repository.Name)
		assert.Equal(t, "tmtk75", w.Repository.Owner.Name)
	})

	t.Run("WebhookProvider", func(t *testing.T) {
		L := "https://gitlab.com/tmtk75"
		R := "google-sheets-sample-go"
		assert.Equal(t, L+"/"+R+"", b.RepositoryURL(w))
		assert.Equal(t, L+"/"+R+"/commits/8c53a97bc3491c125baaacb4974aab5a0ca62c8b", b.CommitURL(w))
		assert.Equal(t, L+"/"+R+"/compare/0000000000000000000000000000000000000000...8c53a97bc3491c125baaacb4974aab5a0ca62c8b", b.CompareURL(w))
		assert.Equal(t, L+"/"+R+"/tree/0.0.8", b.RefURL(w))
		assert.Equal(t, L, b.PusherURL(w))
	})
}
