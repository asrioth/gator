package config

import (
	"encoding/json"
	"os"
)

const gatorConfigName string = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (C *Config) SetUser(userName string) error {
	gatorConfigPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configFile, err := os.Create(gatorConfigPath)
	if err != nil {
		return err
	}
	defer configFile.Close()
	C.CurrentUserName = userName
	encoder := json.NewEncoder(configFile)
	if err := encoder.Encode(C); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	gatorConfigPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configFile, err := os.Open(gatorConfigPath)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	gatorConfigPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	gatorConfigPath += "/" + gatorConfigName
	return gatorConfigPath, nil
}
