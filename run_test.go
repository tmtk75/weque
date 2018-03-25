package weque

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	env := []string{
		"REPOSITORY_NAME=a",
		"OWNER_NAME=b",
		"EVENT=c",
		"DELIVERY=d",
		"REF=e",
		"AFTER=f",
		"BEFORE=g",
		"CREATED=h",
		"DELETED=i",
		"PUSHER_NAME=j",
	}

	expected := []string{
		"REPOSITORY_NAME: a",
		"OWNER_NAME:      b",
		"EVENT:           c",
		"DELIVERY:        d",
		"REF:             e",
		"AFTER:           f",
		"BEFORE:          g",
		"CREATED:         h",
		"DELETED:         i",
		"PUSHER_NAME:     j",
	}

	out := bytes.NewBufferString("")
	stdout = out
	err := Run(env, ".", "./print-env.sh")
	assert.Nil(t, err)

	lines := strings.Split(strings.Trim(out.String(), "\n"), "\n")
	for i, v := range lines {
		assert.Equal(t, expected[i], v)
	}
}
