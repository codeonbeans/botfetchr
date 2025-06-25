# BotVideoSaver

A high-performance Go application that downloads videos from various social media platforms using a browser pool architecture for concurrent task processing.

## 🚀 Features

### Main Features

- **Video Download Bot**: Send a link and get the video downloaded and sent back
- **Browser Pool**: Efficient concurrent task handling with multiple browser instances
- **Proxy Support**: Built-in proxy configuration for each browser instance

### Supported Platforms

Public contents only:

- **Instagram**: Reels, Posts, ~~Stories~~
- **VK**: Videos

## 🏗️ Architecture

### Core Components

#### Browser Pool

The main core of the application is a sophisticated browser pool system that manages multiple Chrome/Chromium instances to handle concurrent video download tasks efficiently.

**Key Features:**

- Round-robin browser selection for load balancing
- Task queuing system with configurable buffer sizes
- Graceful shutdown handling with signal management

#### Browser Management

Each browser instance in the pool:

- Runs in headless or headed mode (configurable)
- Supports proxy authentication
- Handles certificate errors for MITM proxies
- Manages individual task execution with context cancellation

## 📋 Requirements

### System Requirements

Chrome or Chromium browser must be installed and available in your system PATH.

#### Supported Chrome/Chromium Paths

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

## 🛡️ Error Handling

The application includes comprehensive error handling:

- **Panic Recovery**: Tasks that panic are caught and handled gracefully
- **Context Cancellation**: Proper cleanup on shutdown signals
- **Browser Failure**: Individual browser failures don't affect the entire pool
- **Timeout Management**: Configurable timeouts for long-running tasks
