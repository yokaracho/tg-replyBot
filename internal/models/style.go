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
		Name:        "casual",
		DisplayName: "Обычный ответ",
		Description: "Простой человеческий ответ",
		Emoji:       "💬",
		Prompt:      "Ответь на это сообщение как обычный человек в переписке. Будь естественным, используй разговорную речь. Не копируй исходное сообщение - дай свой ответ.",
	},
	{
		Name:        "brief",
		DisplayName: "Краткий ответ",
		Description: "Короткий ответ по существу",
		Emoji:       "⚡",
		Prompt:      "Ответь на это сообщение очень кратко - максимум одно предложение. Будь по существу, но вежливым.",
	},
	{
		Name:        "friendly",
		DisplayName: "Дружелюбный ответ",
		Description: "Теплый дружеский ответ",
		Emoji:       "😊",
		Prompt:      "Ответь дружелюбно и тепло, как близкому другу. Покажи заинтересованность и позитивные эмоции.",
	},
	{
		Name:        "funny",
		DisplayName: "С юмором",
		Description: "Ответ с шуткой или юмором",
		Emoji:       "😄",
		Prompt:      "Ответь с юмором или легкой шуткой. Будь остроумным, но не обидным. Подними настроение собеседнику.",
	},
	{
		Name:        "question",
		DisplayName: "Вопрос",
		Description: "Ответ в виде вопроса",
		Emoji:       "❓",
		Prompt:      "Ответь встречным вопросом или несколькими вопросами. Покажи интерес к теме и желание узнать больше.",
	},
	{
		Name:        "supportive",
		DisplayName: "Поддерживающий",
		Description: "Ответ с пониманием и поддержкой",
		Emoji:       "🤝",
		Prompt:      "Ответь с пониманием и поддержкой. Покажи эмпатию, подбодри человека или просто дай понять, что ты его понимаешь.",
	},
	{
		Name:        "thoughtful",
		DisplayName: "Вдумчивый ответ",
		Description: "Размышляющий содержательный ответ",
		Emoji:       "🤔",
		Prompt:      "Дай вдумчивый ответ с размышлениями по теме. Поделись своим мнением, опытом или соображениями.",
	},
}
