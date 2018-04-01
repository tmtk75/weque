package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/registry"
)

func TestRegistryEnv(t *testing.T) {
	e := &registry.Event{
		ID: "115d94fb-7318-497c-aa74-4061d4e52548",
		Target: registry.Target{
			Repository: "alpine",
			Digest:     "sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93",
			URL:        "http://localhost:5000/v2/alpine/manifests/sha256:d6eda1410b93902ac84bdd775167c84ab59e5abadad88791d742fea93b161e93",
			Tag:        "3.6",
		},
		Request: registry.Request{
			ID:        "329f13c1-1b15-467d-b917-7a79f2113ce3",
			Addr:      "172.17.0.1:54070",
			Host:      "localhost:5000",
			UserAgent: "docker/17.12.0-ce go/go1.9.2 git-commit/c97c6d6 kernel/4.9.60-linuxkit-aufs os/linux arch/amd64 UpstreamClient(Docker-Client/17.12.0-ce \\(darwin\\))",
		},
	}

	s := RegistryEnv(e)

	assert.Equal(t, "EVENT_ID="+e.ID, s[0])
	assert.Equal(t, "REPOSITORY="+e.Target.Repository, s[1])
	assert.Equal(t, "DIGEST="+e.Target.Digest, s[2])
	assert.Equal(t, "URL="+e.Target.URL, s[3])
	assert.Equal(t, "TAG="+e.Target.Tag, s[4])
	assert.Equal(t, "REQUEST_ID="+e.Request.ID, s[5])
	assert.Equal(t, "ADDR="+e.Request.Addr, s[6])
	assert.Equal(t, "HOST="+e.Request.Host, s[7])
	assert.Equal(t, "USER_AGENT="+e.Request.UserAgent, s[8])
}
