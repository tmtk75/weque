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

func TestEventEnv(t *testing.T) {
	e := &Event{
		ID:        "115d94fb-7318-497c-aa74-4061d4e52548",
		Timestamp: "2018-04-01T00:20:31.648394159Z",
		Action:    "push",
		Target: Target{
			Repository: "alpine",
			Digest:     "sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93",
			URL:        "http://localhost:5000/v2/alpine/manifests/sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93",
			Tag:        "3.6",
		},
		Request: Request{
			ID:        "329f13c1-1b15-467d-b917-7a79f2113ce3",
			Addr:      "172.17.0.1:54070",
			Host:      "localhost:5000",
			Method:    "PUT",
			UserAgent: "docker/17.12.0-ce go/go1.9.2 git-commit/c97c6d6 kernel/4.9.60-linuxkit-aufs os/linux arch/amd64 UpstreamClient(Docker-Client/17.12.0-ce \\(darwin\\))",
		},
	}

	s := e.Env()

	assert.Equal(t, "EVENT_ID="+e.ID, s[0])
	assert.Equal(t, "TIMESTAMP="+e.Timestamp, s[1])
	assert.Equal(t, "ACTION="+e.Action, s[2])
	assert.Equal(t, "REPOSITORY="+e.Target.Repository, s[3])
	assert.Equal(t, "DIGEST="+e.Target.Digest, s[4])
	assert.Equal(t, "URL="+e.Target.URL, s[5])
	assert.Equal(t, "TAG="+e.Target.Tag, s[6])
	assert.Equal(t, "REQUEST_ID="+e.Request.ID, s[7])
	assert.Equal(t, "ADDR="+e.Request.Addr, s[8])
	assert.Equal(t, "HOST="+e.Request.Host, s[9])
	assert.Equal(t, "METHOD="+e.Request.Method, s[10])
	assert.Equal(t, "USER_AGENT="+e.Request.UserAgent, s[11])
}
