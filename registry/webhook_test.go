package registry

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	b, err := ioutil.ReadFile("./payload.json")
	assert.NoError(t, err)

	wh, err := Unmarshal(b)

	target := wh.Events[0].Target
	assert.Equal(t, "alpine", target.Repository)
	assert.Equal(t, "3.6", target.Tag)
	assert.Equal(t, "sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93", target.Digest)
	assert.Equal(t, "http://localhost:5000/v2/alpine/manifests/sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93", target.URL)

	req := wh.Events[0].Request
	assert.Equal(t, "329f13c1-1b15-467d-b917-7a79f2113ce3", req.ID)
	assert.Equal(t, "172.17.0.1:54070", req.Addr)
	assert.Equal(t, "localhost:5000", req.Host)
	assert.Equal(t, "docker/17.12.0-ce go/go1.9.2 git-commit/c97c6d6 kernel/4.9.60-linuxkit-aufs os/linux arch/amd64 UpstreamClient(Docker-Client/17.12.0-ce \\(darwin\\))", req.UserAgent)
}
