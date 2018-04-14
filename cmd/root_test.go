package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	assert.Equal(t, ":9981", viper.GetString("port"))
	assert.Equal(t, "", viper.GetString("prefix"))
}
