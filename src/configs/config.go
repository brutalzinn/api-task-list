package configs

import "github.com/spf13/viper"

var cfg *config

type config struct {
	API            APIConfig
	Authentication AuthConfig
}

type APIConfig struct {
	Port       string
	MaxApiKeys int64
}

type AuthConfig struct {
	Secret     string
	Expiration int64
	AesKey     string
}

func init() {
	viper.SetDefault("api.port", 9000)
}

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	cfg = new(config)
	cfg.API = APIConfig{
		Port:       viper.GetString("api.port"),
		MaxApiKeys: viper.GetInt64("api.maxapikeys"),
	}
	cfg.Authentication = AuthConfig{
		Secret:     viper.GetString("authentication.secret"),
		Expiration: viper.GetInt64("authentication.expiration"),
		AesKey:     viper.GetString("authentication.aeskey"),
	}
	return nil
}

func GetConfig() *config {
	return cfg
}

func GetAuthSecret() []byte {
	return []byte(cfg.Authentication.Secret)
}
func GetAesSecret() []byte {
	return []byte(cfg.Authentication.AesKey)
}
