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
	Stdout = out
	err := Run(env, ".", "./bin/print-env.sh")
	assert.Nil(t, err)

	lines := strings.Split(strings.Trim(out.String(), "\n"), "\n")
	for i, v := range lines {
		assert.Equal(t, expected[i], v)
	}
}

func TestRunFailed(t *testing.T) {
	t.Run("missing-script", func(t *testing.T) {
		err := Run([]string{}, ".", "missing-script")
		assert.Error(t, err)
	})

	t.Run("failed-to-run", func(t *testing.T) {
		err := Run([]string{}, ".", "./bin/fail.sh")
		assert.Error(t, err)
		assert.Regexp(t, "^failed to run.*exit status 127", err.Error())
	})
}
