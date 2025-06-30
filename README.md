# botmediasaver

A high-performance Go application that downloads videos from various social media platforms using a browser pool architecture for concurrent task processing. Built with Telegram bot integration for easy user interaction.

## 🚀 Features

### Main Features

-   **Telegram Bot Integration**: Send a link to the bot and get the video downloaded and sent back
-   **Browser Pool Architecture**: Efficient concurrent task handling with multiple browser instances
-   **Multi-Platform Support**: Download videos from Instagram, VK, and more
-   **Proxy Support**: Built-in proxy configuration for both Telegram API and browser instances
-   **High Performance**: Redis-based task queuing and PostgreSQL for data persistence
-   **Smart User Agent Rotation**: Randomized or custom user agents to avoid detection

### Supported Platforms

Public content only:

-   **Instagram**: Reels, Posts, ~~Stories~~
-   **VK**: Videos

## 🏗️ Architecture

### Core Components

#### Telegram Bot

-   Token-based authentication
-   Optional proxy support for API requests
-   Debug logging capabilities
-   SOCKS5 and MTProxy support

#### Browser Pool

The main core of the application is a sophisticated browser pool system that manages multiple Chrome/Chromium instances to handle concurrent video download tasks efficiently.

**Key Features:**

-   Configurable pool size for concurrent browser instances
-   Task queuing system with configurable buffer sizes
-   Headless or headed mode operation
-   Per-instance proxy configuration
-   Graceful shutdown handling with signal management

#### Media Saver Engine

-   Smart user agent rotation (random or predefined list)
-   Quality selection (low/high)
-   Configurable retry mechanism
-   Timeout management for reliable downloads

#### Data Layer

-   **PostgreSQL**: Persistent data storage with connection pooling
-   **Redis**: Task queuing and caching with cluster support
-   **Migrations**: Versioned database schema management

## 🚀 Quick Start

### Prerequisites

Before running the application, you need to create a configuration file:

1.  Copy the example configuration:

    ```bash
    cp config.example.yml config.dev.yml

    ```

2.  Edit `config.dev.yml` with your settings (see Configuration section below)

### Option 1: Docker Compose (Recommended)

The easiest way to run botmediasaver:

```bash
docker compose -f docker-compose.yml -p botmediasaver up -d

```

**Advantages:**

-   No system dependencies required
-   All tools pre-installed (Chrome, PostgreSQL, Redis)
-   Easy to manage and scale

### Option 2: Native Installation

For development or custom deployments:

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/botmediasaver.git
cd botmediasaver

# 2. Install dependencies
go mod download

# 3. Install required tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 4. Set up database
cd migrations
goose postgres "your-postgres-connection-string" up
cd ..
sqlc generate

# 5. Build and run
go build -o botmediasaver
./botmediasaver

```

## 📋 Requirements

### Docker Requirements

When using Docker Compose:

-   Docker and Docker Compose installed
-   Configuration file (`config.dev.yml`)
-   Telegram Bot Token

**That's it!** The container includes all necessary dependencies.

### Native Installation Requirements

When running natively, you need:

#### 1. External Services

-   **PostgreSQL** database server
-   **Redis** server (single instance)
-   **Telegram Bot Token** (from @BotFather)

#### 2. Go Tools

-   **Goose** - Database migration tool
-   **SQLC** - Generate type-safe Go code from SQL

#### 3. Chrome/Chromium Browser

Must be installed and available in your system PATH.

<details> <summary>📁 Supported Browser Paths</summary>

**macOS:**

```
/Applications/Google Chrome.app/Contents/MacOS/Google Chrome
/Applications/Chromium.app/Contents/MacOS/Chromium
/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge
/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary
/usr/bin/google-chrome-stable
/usr/bin/google-chrome
/usr/bin/chromium
/usr/bin/chromium-browser

```

**Linux:**

```
chrome
google-chrome
/usr/bin/google-chrome
microsoft-edge
/usr/bin/microsoft-edge
chromium
chromium-browser
/usr/bin/google-chrome-stable
/usr/bin/chromium
/usr/bin/chromium-browser
/snap/bin/chromium
/data/data/com.termux/files/usr/bin/chromium-browser

```

**Windows:**

```
chrome
edge
%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe
%LOCALAPPDATA%\Chromium\Application\chrome.exe
%PROGRAMFILES%\Google\Chrome\Application\chrome.exe
%PROGRAMFILES(X86)%\Google\Chrome\Application\chrome.exe
%LOCALAPPDATA%\Microsoft\Edge\Application\msedge.exe
%PROGRAMFILES%\Microsoft\Edge\Application\msedge.exe
%PROGRAMFILES(X86)%\Microsoft\Edge\Application\msedge.exe

```

**OpenBSD:**

```
chrome
chromium

```

</details>

## ⚙️ Configuration

The application is configured via YAML files. Copy `config.example.yml` to your environment-specific config file and customize the following sections:

### Environment Settings

```yaml
env: "dev" # Available options: dev, staging, production
app:
  name: "botmediasaver"

```

### Telegram Bot Configuration

```yaml
telegramBot:
  token: "<your-telegram-bot-token>" # Get from @BotFather
  logDebug: false
  proxy: # Optional proxy settings
    enabled: false
    type: "socks5" # socks5 or mtproxy
    address: ""
    port: 1080
    username: "" # For authenticated proxies
    password: ""

```

### Media Saver Settings

```yaml
mediaSaver:
  useRandomUA: true # Random user agent per request
  userAgents: [] # Custom user agents (fallback to random if empty)
  quality: "high" # low or high
  retryCount: 3 # Failed task retries
  timeout: 15 # Seconds

```

### Database Configuration

```yaml
postgres:
  # Option 1: Connection URL
  url: "postgresql://username:password@localhost:5432/dbname"

  # Option 2: Individual parameters
  host: "localhost" # Use "host.docker.internal" for Docker
  port: "5432"
  database: "botmediasaver"
  username: "your_username"
  password: "your_password"

  # Connection pooling
  maxConnections: 8
  maxIdleConnections: 10
  logQuery: false # Enable for debugging

```

### Redis Configuration

```yaml
redis:
  clusters:
    - host: "localhost" # Use "host.docker.internal" for Docker
      port: "6379"
  password: ""
  db: 0

```

### Browser Pool Settings

```yaml
browserpool:
  headless: true # Set false to see browser UI
  poolSize: 10 # Concurrent browser instances
  proxies: [] # List of proxy URLs (≤ poolSize)
  taskQueueSize: 5 # Tasks per browser instance

```

### Logging Configuration

```yaml
log:
  level: "debug" # debug, info, warn, error, dpanic, panic, fatal
  stacktraceLevel: "error"
  fileEnabled: true
  fileSize: 10 # MB
  filePath: "log/log.log"
  fileCompress: false
  maxAge: 1 # Days to keep log files
  maxBackups: 10 # Number of log files

```

## 🔧 Development

### Database Management

#### Creating Migrations

```bash
cd migrations
goose create add_new_table sql

```

#### Applying Migrations

```bash
# From migrations directory
goose postgres "your-connection-string" up

# Rollback if needed
goose postgres "your-connection-string" down

```

#### Generating Code

After modifying SQL queries:

```bash
sqlc generate

```

### Environment-Specific Configs

Create different config files for different environments:

-   `config.dev.yml` - Development
-   `config.staging.yml` - Staging
-   `config.prod.yml` - Production

Set the environment via the `env` field in your config file.

## 🛡️ Error Handling

The application includes comprehensive error handling:

-   **Panic Recovery**: Tasks that panic are caught and handled gracefully
-   **Context Cancellation**: Proper cleanup on shutdown signals
-   **Browser Failure**: Individual browser failures don't affect the entire pool
-   **Timeout Management**: Configurable timeouts for long-running tasks
-   **Retry Logic**: Automatic retries for failed download attempts
-   **Connection Pooling**: Database connection management with automatic recovery

## 🚀 Scaling & Performance

### Horizontal Scaling

-   Increase `browserpool.poolSize` for more concurrent downloads
-   Add more Redis cluster nodes for better task distribution
-   Use database connection pooling for optimal resource usage

### Performance Tuning

-   Adjust `taskQueueSize` based on memory constraints
-   Optimize `retryCount` and `timeout` values for your use case
-   Use headless mode (`browserpool.headless: true`) for better performance

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://claude.ai/chat/LICENSE) file for details.
