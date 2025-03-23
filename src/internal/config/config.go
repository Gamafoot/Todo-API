package config

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Http struct {
		Port           string        `yaml:"port" env-default:"8000"`
		ReadTimeout    time.Duration `yaml:"read_timeout" env-default:"10s"`
		WriteTimeout   time.Duration `yaml:"write_timeout" env-default:"10s"`
		MaxHeaderBytes int           `yaml:"max_header_bytes" env-default:"1"`
	} `yaml:"http"`

	Database struct {
		URL string `env:"url" env-required:"true"`
	} `yaml:"database"`

	Hash struct {
		Salt string `yaml:"salt" env-required:"true"`
	} `yaml:"hash"`

	Jwt struct {
		SigningKey      string        `yaml:"signing_key" env-required:"true"`
		AccessTokenTtl  time.Duration `yaml:"access_token_ttl" env-required:"true"`
		RefreshTokenTtl time.Duration `yaml:"refresh_token_ttl" env-required:"true"`
	} `yaml:"jwt"`

	Cors struct {
		AllowOrigins []string `yaml:"origins"`
		AllowMethods []string `yaml:"methods"`
		AllowHeaders []string `yaml:"headers"`
	} `yaml:"cors"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		path := getConfigPath()

		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "../config/default.yaml", "set config file")

	envPath := os.Getenv("CONFIG_PATH")

	if len(envPath) > 0 {
		path = envPath
	}

	return path
}
