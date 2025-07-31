package services

import (
	"fmt"

	"tg-replyBot/internal/models"
)

type StyleManager struct {
	styles map[string]models.Style
}

func NewStyleManager() *StyleManager {
	styles := make(map[string]models.Style)
	for _, style := range models.DefaultStyles {
		styles[style.Name] = style
	}

	return &StyleManager{
		styles: styles,
	}
}

func (sm *StyleManager) GetAllStyles() []models.Style {
	styles := make([]models.Style, 0, len(sm.styles))
	for _, style := range models.DefaultStyles {
		styles = append(styles, style)
	}
	return styles
}

func (sm *StyleManager) GetStyle(name string) (models.Style, error) {
	style, exists := sm.styles[name]
	if !exists {
		return models.Style{}, fmt.Errorf("стиль '%s' не найден", name)
	}
	return style, nil
}

func (sm *StyleManager) GetMainStyles(count int) []models.Style {
	allStyles := sm.GetAllStyles()
	if len(allStyles) <= count {
		return allStyles
	}
	return allStyles[:count]
}
