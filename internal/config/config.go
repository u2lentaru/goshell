package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env   string `yaml:"env" env-default:"development"`
	PgURL string `yaml:"pg_url" env-required:"true"`
	// HTTPServer  `yaml:"http_server"`
}

// type HTTPServer struct {
// 	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
// 	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
// 	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
// 	User        string        `yaml:"user" env-required:"true"`
// 	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
// }

func MustLoad() *Config {
	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	// configPath := os.Getenv("CONFIG_PATH")
	configPath := "../../config/local/yaml"

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
