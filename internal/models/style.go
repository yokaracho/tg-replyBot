package models

type Style struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	Prompt      string `json:"prompt"`
}

var DefaultStyles = []Style{
	{
		Name:        "variant1",
		DisplayName: "Вариант 1",
		Description: "Скажи это по-своему, сохраняя смысл",
		Emoji:       "🔹",
		Prompt:      "Передай смысл сообщения своими словами. Формулировка может быть другой, но идея должна остаться той же.",
	},
	{
		Name:        "variant2",
		DisplayName: "Вариант 2",
		Description: "Переиначь сообщение по-другому",
		Emoji:       "🔸",
		Prompt:      "Переиначь это сообщение по-другому, сохранив основной смысл. Избегай копирования оригинального текста.",
	},
	{
		Name:        "variant3",
		DisplayName: "Вариант 3",
		Description: "Вырази ту же мысль по-другому",
		Emoji:       "✴️",
		Prompt:      "Вырази ту же самую мысль, но по-другому — используй другие слова, стиль и структуру.",
	},
	{
		Name:        "variant4",
		DisplayName: "Вариант 4",
		Description: "Свободный пересказ",
		Emoji:       "🔁",
		Prompt:      "Сделай пересказ сообщения в свободной форме, сохраняя интонацию и отношение, но не копируя фразы.",
	},
	{
		Name:        "variant5",
		DisplayName: "Вариант 5",
		Description: "Переоформление мысли",
		Emoji:       "🌀",
		Prompt:      "Переоформи основную мысль сообщения, добавив немного от себя, но не уходя от сути.",
	},
}
