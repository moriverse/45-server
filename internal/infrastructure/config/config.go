package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Log      LogConfig
	Wechat   WechatConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	SecretKey      string `mapstructure:"secret_key"`
	ExpiresInHours int    `mapstructure:"expires_in_hours"`
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type LogConfig struct {
	Level  string
	Format string
}

type WechatConfig struct {
	AppID            string `mapstructure:"app_id"`
	AppSecret        string `mapstructure:"app_secret"`
	Host             string `mapstructure:"host"`
	Code2SessionPath string `mapstructure:"code2session_path"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
