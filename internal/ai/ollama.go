package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tg-replyBot/internal/config"
	"tg-replyBot/pkg/logger"
)

type Ollama struct {
	config config.OllamaConfig
	client *http.Client
	logger logger.Logger
}

type ollamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func NewOllama(config config.OllamaConfig, logger logger.Logger) Provider {
	return &Ollama{
		config: config,
		client: &http.Client{Timeout: 60 * time.Second}, // Увеличиваем таймаут для локальных моделей
		logger: logger,
	}
}

func (o *Ollama) GenerateReply(ctx context.Context, request Request) (string, error) {
	prompt := o.buildPrompt(request)

	reqBody := ollamaRequest{
		Model:  o.config.Model,
		Prompt: prompt,
		Stream: false, // Отключаем стриминг для простоты
		Options: map[string]interface{}{
			"temperature": o.config.Temperature,
			"num_predict": o.config.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		o.logger.Error("Ошибка сериализации запроса к Ollama", "error", err)
		return "", fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		o.logger.Error("Ошибка создания запроса к Ollama", "error", err)
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	o.logger.Debug("Отправляем запрос к Ollama", "url", url, "model", o.config.Model)

	resp, err := o.client.Do(req)
	if err != nil {
		o.logger.Error("Ошибка выполнения запроса к Ollama", "error", err)
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.logger.Error("Ollama вернул ошибку", "status", resp.StatusCode)
		return "", fmt.Errorf("Ollama API вернул статус: %d", resp.StatusCode)
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		o.logger.Error("Ошибка декодирования ответа от Ollama", "error", err)
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if ollamaResp.Error != "" {
		o.logger.Error("Ошибка API Ollama", "error", ollamaResp.Error)
		return "", fmt.Errorf("ошибка API Ollama: %s", ollamaResp.Error)
	}

	if ollamaResp.Response == "" {
		o.logger.Warn("Ollama вернул пустой ответ")
		return "", fmt.Errorf("получен пустой ответ от Ollama")
	}

	response := strings.TrimSpace(ollamaResp.Response)
	o.logger.Debug("Получен ответ от Ollama", "response_length", len(response))

	return response, nil
}

func (o *Ollama) buildPrompt(request Request) string {
	var contextStr string
	if len(request.ContextMessages) > 0 {
		contextStr = "Контекст предыдущих сообщений:\n" + strings.Join(request.ContextMessages, "\n") + "\n\n"
	}

	return fmt.Sprintf(`%sНа следующее сообщение нужно ответить:
"%s"

%s. Ответ должен быть уместным, тактичным и помогать поддержать разговор.
Ответь только текстом ответа, без дополнительных пояснений.`,
		contextStr, request.Message, request.Style.Prompt)
}
