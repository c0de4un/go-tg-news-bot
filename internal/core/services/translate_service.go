package services

import (
	"bufio"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"log"
	"os"
	"strings"
	"sync"
)

type TranslateService struct {
	Dictionaries map[string]*models.Dictionary
}

var (
	translateServiceOnce     sync.Once
	translateServiceInstance *TranslateService
)

func InitializeTranslateService() {
	translateServiceOnce.Do(func() {
		translateServiceInstance = &TranslateService{
			Dictionaries: make(map[string]*models.Dictionary),
		}

		translateServiceInstance.loadDictionaries()
	})
}

func Translate(text string, lng ...string) string {
	language := "ru"
	if len(lng) > 0 {
		language = lng[0]
	}

	return translateServiceInstance.translate(text, language)
}

func (ts *TranslateService) translate(text string, lng string) string {
	dict, exists := ts.Dictionaries[lng]
	if !exists {
		return text
	}

	if translated, ok := dict.Data[text]; ok {
		return translated
	}
	return text
}

func (ts *TranslateService) loadDictionaries() {
	ruDict := &models.Dictionary{
		Language: "ru", // Target language code
		Data:     make(map[string]string),
	}

	file, err := os.Open("config/locale/ru.txt")
	if err != nil {
		log.Fatalf("Failed to open translations file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	var key string

	for scanner.Scan() {
		// Trim quotes and whitespace from line
		line := strings.Trim(scanner.Text(), `" `)

		if line == "" {
			continue // Skip empty lines
		}

		if key == "" {
			key = line // First line of pair (source text)
		} else {
			// Second line of pair (translated text)
			ruDict.Data[key] = line
			key = "" // Reset for next pair
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading translations file: %v", err)
	}

	// Add the loaded dictionary to the service
	ts.Dictionaries[ruDict.Language] = ruDict
}
