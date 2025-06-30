
# botmediasaver

A high-performance Go application that downloads videos from various social media platforms using a browser pool architecture for concurrent task processing.

## 🚀 Features

### Main Features

-   **Video Download Bot**: Send a link and get the video downloaded and sent back
-   **Browser Pool**: Efficient concurrent task handling with multiple browser instances
-   **Proxy Support**: Built-in proxy configuration for each browser instance

### Supported Platforms

Public contents only:

-   **Instagram**: Reels, Posts, ~~Stories~~
-   **VK**: Videos

## 🏗️ Architecture

### Core Components

#### Browser Pool

The main core of the application is a sophisticated browser pool system that manages multiple Chrome/Chromium instances to handle concurrent video download tasks efficiently.

**Key Features:**

-   Round-robin browser selection for load balancing
-   Task queuing system with configurable buffer sizes
-   Graceful shutdown handling with signal management

#### Browser Management

Each browser instance in the pool:

-   Runs in headless or headed mode (configurable)
-   Supports proxy authentication
-   Handles certificate errors for MITM proxies
-   Manages individual task execution with context cancellation

## 📋 Requirements

### Native Installation Requirements

If you're running the application natively (without Docker), you'll need:

#### System Requirements

Chrome or Chromium browser must be installed and available in your system PATH.

##### Supported Chrome/Chromium Paths

**macOS (Darwin):**

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

**OpenBSD:**

```
chrome
chromium

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

#### Development Tools

You'll also need to install these Go tools for database management:

1.  **Goose** - Database migration tool

    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest

    ```

2.  **SQLC** - Generate type-safe Go code from SQL

    ```bash
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

    ```


### Docker Installation

If you're using Docker, **none of the above requirements are needed**. The Docker container includes all necessary dependencies including Chrome/Chromium and the required Go tools.

## 🚀 Usage

There are two ways to run this bot but always to create the configuration first: Create a configuration file **config.dev.yml** base on config/config.example.yml


### Using Docker Compose (Recommended)

The easiest way to run botmediasaver is using Docker Compose:
```bash
# Run the docker compose
docker compose -f docker-compose.yml -p botmediasaver up -d
```

### Native Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/botmediasaver.git
cd botmediasaver
```

#### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install required tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

```

#### 3. Database Setup

```bash
# Run database migrations
goose up

# Generate database code (if you make changes to SQL files)
sqlc generate

```

#### 4. Build and Run

```bash
# Build the application
go build -o botmediasaver

# Run the application
./botmediasaver

```

### Development Workflow

#### Making Database Changes

1.  Create a new migration:

    ```bash
    cd migrations
    goose create add_new_table sql
    ```

2.  Edit the generated migration file in `migrations/`

3.  Apply migrations:

    ```bash
    goose up
    ```

4.  If you modify SQL queries, regenerate the Go code:

    ```bash
    sqlc generate
    ```


## 🛡️ Error Handling

The application includes comprehensive error handling:

-   **Panic Recovery**: Tasks that panic are caught and handled gracefully
-   **Context Cancellation**: Proper cleanup on shutdown signals
-   **Browser Failure**: Individual browser failures don't affect the entire pool
-   **Timeout Management**: Configurable timeouts for long-running tasks
