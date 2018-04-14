package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	os.Setenv("INSECURE", "true")
	assert.True(t, viper.GetBool("INSECURE"))

	os.Setenv("INSECURE", "false")
	assert.False(t, viper.GetBool("INSECURE"))
}
