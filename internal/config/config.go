package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	AI       AIConfig       `yaml:"ai"`
	Ollama   OllamaConfig   `yaml:"ollama"`
	Storage  StorageConfig  `yaml:"storage"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type TelegramConfig struct {
	Token   string `yaml:"token" env:"TELEGRAM_BOT_TOKEN"`
	Debug   bool   `yaml:"debug" env:"TELEGRAM_DEBUG" default:"false"`
	Timeout int    `yaml:"timeout" env:"TELEGRAM_TIMEOUT" default:"60"`
}

type AIConfig struct {
	PrimaryProvider string `yaml:"primary_provider" env:"AI_PRIMARY_PROVIDER" default:"ollama"`
	SmartFallback   bool   `yaml:"smart_fallback" env:"AI_SMART_FALLBACK" default:"true"`
}

type OllamaConfig struct {
	BaseURL     string  `yaml:"base_url" env:"OLLAMA_BASE_URL" default:"http://localhost:11434"`
	Model       string  `yaml:"model" env:"OLLAMA_MODEL" default:"llama2"`
	Temperature float64 `yaml:"temperature" env:"OLLAMA_TEMPERATURE" default:"0.7"`
	MaxTokens   int     `yaml:"max_tokens" env:"OLLAMA_MAX_TOKENS" default:"800"`
	Timeout     int     `yaml:"timeout" env:"OLLAMA_TIMEOUT" default:"60"`
}

type StorageConfig struct {
	Type   string       `yaml:"type" env:"STORAGE_TYPE" default:"memory"`
	Memory MemoryConfig `yaml:"memory"`
	Redis  RedisConfig  `yaml:"redis"`
}

type MemoryConfig struct {
	CleanupInterval time.Duration `yaml:"cleanup_interval" default:"30m"`
	TTL             time.Duration `yaml:"ttl" default:"24h"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" default:"localhost"`
	Port     int    `yaml:"port" env:"REDIS_PORT" default:"6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" default:"0"`
}

type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" default:"json"`
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{
		Telegram: TelegramConfig{
			Debug:   false,
			Timeout: 60,
		},
		Ollama: OllamaConfig{
			BaseURL:     "http://localhost:11434",
			Model:       "llama3",
			MaxTokens:   800,
			Temperature: 0.7,
		},
		Storage: StorageConfig{
			Type: "memory",
			Memory: MemoryConfig{
				CleanupInterval: 30 * time.Minute,
				TTL:             24 * time.Hour,
			},
		},
		Logger: LoggerConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Загрузка из YAML-файла, если указан
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		data, err := os.ReadFile(configFile)
		if err == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	}

	// Переопределения из переменных окружения
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		cfg.Telegram.Token = token
	}

	if debug := os.Getenv("TELEGRAM_DEBUG"); debug != "" {
		if parsed, err := strconv.ParseBool(debug); err == nil {
			cfg.Telegram.Debug = parsed
		}
	}

	if baseURL := os.Getenv("OLLAMA_BASE_URL"); baseURL != "" {
		cfg.Ollama.BaseURL = baseURL
	}

	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		cfg.Ollama.Model = model
	}

	return cfg, nil
}
