package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tg-replyBot/internal/ai"
	"tg-replyBot/internal/bot"
	"tg-replyBot/internal/config"
	"tg-replyBot/internal/services"
	"tg-replyBot/internal/storage/memory"
	"tg-replyBot/pkg/logger"
)

func main() {
	printOutboundIPInfo()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	logger := logger.New(cfg.Logger)

	// Инициализация хранилища
	storage := memory.New()

	// Инициализация AI провайдера (только Ollama)
	aiProvider := ai.NewOllama(cfg.Ollama, logger)

	// Инициализация сервисов
	contextManager := services.NewContextManager(storage, logger)
	styleManager := services.NewStyleManager()
	replyGenerator := services.NewReplyGenerator(aiProvider, styleManager, logger)

	// Создание бота
	telegramBot, err := bot.New(cfg.Telegram, contextManager, replyGenerator, logger)
	if err != nil {
		logger.Fatal("Ошибка создания бота", "error", err)
	}

	// Запуск бота в горутине
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := telegramBot.Start(ctx); err != nil {
			logger.Error("Ошибка работы бота", "error", err)
		}
	}()

	logger.Info("Бот запущен")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Остановка бота...")
	cancel()
	telegramBot.Stop()
	logger.Info("Бот остановлен")
}

type IPInfo struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Org     string `json:"org"`
}

func printOutboundIPInfo() {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		fmt.Println("Ошибка запроса к ipinfo.io:", err)
		return
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}

	fmt.Println("📡 IP-адрес и страна, видимые снаружи:")
	fmt.Printf("IP: %s\n", info.IP)
	fmt.Printf("Страна: %s\n", info.Country)
	fmt.Printf("Регион: %s\n", info.Region)
	fmt.Printf("Город: %s\n", info.City)
	fmt.Printf("Организация: %s\n", info.Org)
}
