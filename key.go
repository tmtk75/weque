package weque

import "github.com/spf13/viper"

const (
	KeySecretToken  = "secret_token"
	KeyInsecureMode = "insecure" // Skip verification step for development
	KeyPrefix       = "prefix"
)

func init() {
	viper.BindEnv(KeySecretToken, "SECRET_TOKEN")
}
