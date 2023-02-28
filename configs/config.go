package configs

import "github.com/spf13/viper"

var cfg *config

type config struct {
	API            APIConfig
	DB             DBConfig
	Authentication AuthConfig
	Redis          RedisConfig
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
type RedisConfig struct {
	Host string
	Db   int
}
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

func init() {
	viper.SetDefault("api.port", 9000)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.db", 1)

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
	cfg.DB = DBConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Pass:     viper.GetString("database.pass"),
		Database: viper.GetString("database.database"),
	}
	cfg.Redis = RedisConfig{
		Host: viper.GetString("redis.host"),
		Db:   viper.GetInt("redis.db"),
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
