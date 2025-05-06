package tg

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type TelegramConfig struct {
	XMLName     xml.Name `xml:"TelegramConfig"`
	EditorToken string   `xml:"EditorToken"`
	EditorBotID int64    `xml:"EditorBotID"`
	ReaderToken string   `xml:"ReaderToken"`
	ReaderBotID int64    `xml:"ReaderBotID"`
	AdminID     int64    `xml:"AdminID"`
}

func LoadTelegramConfig() (*TelegramConfig, error) {
	cfgPath := filepath.Join("config", "tg.xml")
	content, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Printf("failed to load tg config: %s", err)
		return nil, err
	}

	model := &TelegramConfig{}
	if err := xml.Unmarshal(content, &model); err != nil {
		fmt.Printf("failed to parse tg config: %s", err)
		return nil, err
	}

	return model, nil
}
