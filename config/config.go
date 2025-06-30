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
	VideoSaver  VideoSaver  `yaml:"videoSaver"`
	Log         Log         `yaml:"log"`
	Postgres    Postgres    `yaml:"postgres"`
	Redis       Redis       `yaml:"redis"`
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

type VideoSaver struct {
	UseRandomUA       bool     `yaml:"useRandomUA"`       // Use random user agent for
	UserAgents        []string `yaml:"userAgents"`        // List of user agents to use if UseRandomUA is false
	Quality           string   `yaml:"quality"`           // Available options: low, high
	RetryCount        int      `yaml:"retryCount"`        // Number of retries for failed
	Timeout           int      `yaml:"timeout"`           // Timeout in seconds for each download
	MaxGroupMediaSize int64    `yaml:"maxGroupMediaSize"` // Maximum size of media group in MB, if the group exceeds this size, it will be split into multiple messages
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

type Redis struct {
	Host     string `yaml:"host"`     // Redis host, use "host.docker.internal" if you run app inside docker container
	Port     string `yaml:"port"`     // Redis port
	Password string `yaml:"password"` // Redis password
	DB       int    `yaml:"db"`       // Redis database number, default is 0
}

type BrowserPool struct {
	Headless      bool     `yaml:"headless"` // Whether to run browsers in headless mode
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
