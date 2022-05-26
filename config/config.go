package config

import (
	"log"

	"github.com/arrow2nd/nekome/oauth"
)

// Config 設定
type Config struct {
	// Cred 認証情報
	Cred *Cred
	// Settings 設定情報
	Settings *Settings
	dirPath  string
}

// New 生成
func New() *Config {
	path, err := getConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Cred: &Cred{
			tokens: map[string]*oauth.Token{},
		},
		Settings: &Settings{},
		dirPath:  path,
	}
}
