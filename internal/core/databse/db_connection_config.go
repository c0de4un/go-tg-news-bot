package database

import (
	"encoding/xml"
	"fmt"
	"os"
)

type dbConnectionConfig struct {
	XMLName  xml.Name `xml:"Connection"`
	Host     string   `xml:"Host"`
	Port     int      `xml:"Port"`
	User     string   `xml:"User"`
	Password string   `xml:"Password"`
	DBName   string   `xml:"DB"`
}

func loadDBConfig(path string) (*dbConnectionConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error at loadDBConfig, while closing connection: %s", err)
		}
	}(file)

	var config dbConnectionConfig
	if err := xml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %v", err)
	}

	return &config, nil
}
