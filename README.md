# tg-replyBot

Telegram Reply Bot — a simple bot for automatically replying to messages on Telegram.

## Description

This bot allows you to automatically reply to incoming messages on Telegram. It is built with Go using the official Telegram Bot API. Perfect for quickly deploying a simple auto-responder without complex logic.

## Features

- Automatic reply to incoming messages  
- Supports text messages  
- Simple and clear architecture  
- Easy to configure and extend

## Bot Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yokaracho/tg-replyBot.git
   cd tg-replyBot
   ```

2. **Copy the example environment file and add your bot token:**

   ```bash
   cp .env.example .env
   # Edit the .env file and add your TELEGRAM_BOT_TOKEN
   ```

3. **Install dependencies and run the bot using Make:**

   ```bash
   make run
   ```

## Optional: Install Ollama (for local LLM integration)

If you plan to use local large language models (LLMs) with the bot, you can install [Ollama](https://ollama.com) and run a model like `llama2`.

### Installation and Setup

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Start the Ollama server
ollama serve

# Download the LLaMA 3 model
ollama pull llama3

# Run the LLaMa 3 model
ollama run llama3
```
