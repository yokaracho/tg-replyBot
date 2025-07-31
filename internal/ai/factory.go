package ai

import (
	"fmt"

	"tg-replyBot/internal/config"
	"tg-replyBot/pkg/logger"
)

type ProviderType string

const (
	ProviderTypeOllama   ProviderType = "ollama"
	ProviderTypeFallback ProviderType = "fallback"
)

type ProviderFactory struct {
	config config.Config
	logger logger.Logger
}

func NewProviderFactory(cfg config.Config, logger logger.Logger) *ProviderFactory {
	return &ProviderFactory{
		config: cfg,
		logger: logger,
	}
}

func (f *ProviderFactory) CreateProvider(providerType ProviderType) (Provider, error) {
	switch providerType {
	case ProviderTypeOllama:
		return NewOllama(f.config.Ollama, f.logger), nil
	case ProviderTypeFallback:
		return NewFallback(f.logger), nil
	default:
		return nil, fmt.Errorf("неизвестный тип провайдера: %s", providerType)
	}
}

func (f *ProviderFactory) CreatePrimaryProvider() (Provider, error) {
	primaryType := ProviderType(f.config.AI.PrimaryProvider)

	if f.config.AI.SmartFallback && primaryType != ProviderTypeFallback {
		primary, err := f.CreateProvider(primaryType)
		if err != nil {
			f.logger.Warn("Не удалось создать основной провайдер, используем fallback",
				"provider", primaryType, "error", err)
			return f.CreateProvider(ProviderTypeFallback)
		}

		fallback, err := f.CreateProvider(ProviderTypeFallback)
		if err != nil {
			f.logger.Error("Не удалось создать fallback провайдер", "error", err)
			return primary, nil
		}

		f.logger.Info("Создан SmartFallback провайдер",
			"primary", primaryType, "fallback", ProviderTypeFallback)
		return NewSmartFallback(primary, fallback, f.logger), nil
	}

	provider, err := f.CreateProvider(primaryType)
	if err != nil {
		f.logger.Error("Не удалось создать провайдер", "provider", primaryType, "error", err)
		f.logger.Info("Используем fallback провайдер как запасной")
		return f.CreateProvider(ProviderTypeFallback)
	}

	f.logger.Info("Создан провайдер", "type", primaryType)
	return provider, nil
}

func (f *ProviderFactory) ValidateConfig() error {
	primaryType := ProviderType(f.config.AI.PrimaryProvider)

	switch primaryType {
	case ProviderTypeOllama:
		if f.config.Ollama.BaseURL == "" {
			return fmt.Errorf("для Ollama провайдера требуется base_url")
		}
		if f.config.Ollama.Model == "" {
			return fmt.Errorf("для Ollama провайдера требуется model")
		}
	case ProviderTypeFallback:
	default:
		return fmt.Errorf("неподдерживаемый тип провайдера: %s", primaryType)
	}

	return nil
}
