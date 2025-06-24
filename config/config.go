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
	Postgres    Postgres    `yaml:"postgres"`
	BrowserPool BrowserPool `yaml:"browserPool"`
}

type App struct {
	Name string `yaml:"name"`
}

type TelegramBot struct {
	Token    string           `yaml:"token"`
	LogDebug bool             `yaml:"logDebug"`
	Proxy    TelegramBotProxy `yaml:"proxy"`
}

type TelegramBotProxy struct {
	Enabled  bool   `yaml:"enabled"`  // true or false
	Type     string `yaml:"type"`     // "socks5", "http"
	Address  string `yaml:"address"`  //
	Port     int    `yaml:"port"`     // e.g: 1080
	Username string `yaml:"username"` // optional
	Password string `yaml:"password"` // optional
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

type Postgres struct {
	Url             string `yaml:"url"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxConnections  int32  `yaml:"maxConnections"`
	MaxConnIdleTime int32  `yaml:"maxConnIdleTime"`
}

type BrowserPool struct {
	PoolSize      int      `yaml:"poolSize"`
	Proxies       []string `yaml:"proxies"`
	TaskQueueSize int      `yaml:"taskQueueSize"` // Buffer size for task channels
}

func GetConfig() *Config {
	if config == nil {
		SetConfig("config/config.dev.yml")
	}

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
