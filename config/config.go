package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var config *Config
var m sync.Mutex

type Config struct {
	Env         string      `yaml:"env"`
	App         App         `yaml:"app"`
	TelegramBot TelegramBot `yaml:"telegramBot"`
	Log         Log         `yaml:"log"`
}

type App struct {
	Name string `yaml:"name"`
}

type TelegramBot struct {
	Token string `yaml:"token"`
}

type Log struct {
	Level           string `yaml:"level"`
	StacktraceLevel string `yaml:"stacktraceLevel"`
	FileEnabled     bool   `yaml:"fileEnabled"`
	FileSize        int    `yaml:"fileSize"`
	FilePath        string `yaml:"filePath"`
	FileCompress    bool   `yaml:"fileCompress"`
	MaxAge          int    `yaml:"maxAge"`
	MaxBackups      int    `yaml:"maxBackups"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	/** Because GitHub Actions doesn't have .env, and it will load ENV variables from GitHub Secrets */
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
