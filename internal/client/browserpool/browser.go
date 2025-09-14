package browserpool

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/codeonbeans/botfetchr/internal/logger"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Browser struct {
	*rod.Browser
	Proxy string
}

func NewBrowser(headless bool, proxy string) (*Browser, error) {
	chrome, found := launcher.LookPath()
	if !found {
		return nil, fmt.Errorf("could not find Chrome executable in PATH")
	}

	customLauncher := launcher.
		New().
		Bin(chrome).
		Headless(headless)

	if proxy != "" {
		p, err := ParseProxy(proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy: %w", err)
		}

		customLauncher.
			Proxy(fmt.Sprintf("%s:%d", p.Host, p.Port)).
			Delete("use-mock-keychain")
	}

	url, err := customLauncher.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url)

	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	if proxy != "" {
		// we don't have to self issue certs for MITM
		if err = browser.IgnoreCertErrors(true); err != nil {
			return nil, fmt.Errorf("failed to ignore cert errors: %w", err)
		}

		// Adding authentication to the proxy, for the next auth request.
		// We use CLI tool "mitmproxy --proxyauth user:pass" as an example.
		go browser.MustHandleAuth("user", "pass")()
	}

	// if _, err = browser.Page(proto.TargetCreateTarget{
	// 	URL: "about:blank",
	// }); err != nil {
	// 	return nil, fmt.Errorf("failed to create empty page: %w", err)
	// }

	return &Browser{
		Browser: browser,
		Proxy:   proxy,
	}, nil
}

type Proxy struct {
	Protocol string // e.g., "http", "https", "socks5"
	Host     string // e.g., "
	Port     int    // e.g., 8080
	Username string // Optional, for authenticated proxies
	Password string // Optional, for authenticated proxies
}

// protocol://username:password@host:port

func ParseProxy(proxyStr string) (Proxy, error) {
	// Example: "http://user:pass@host:port"
	parts := strings.Split(proxyStr, "://")
	if len(parts) != 2 {
		return Proxy{}, fmt.Errorf("invalid proxy format")
	}

	protocol := parts[0]
	authAndHost := parts[1]

	var username, password, host string
	if strings.Contains(authAndHost, "@") {
		authParts := strings.SplitN(authAndHost, "@", 2)
		auth := authParts[0]
		host = authParts[1]

		if strings.Contains(auth, ":") {
			authSplit := strings.SplitN(auth, ":", 2)
			username = authSplit[0]
			password = authSplit[1]
		} else {
			username = auth
		}
	} else {
		host = authAndHost
	}

	hostParts := strings.Split(host, ":")
	if len(hostParts) != 2 {
		return Proxy{}, fmt.Errorf("invalid proxy host format")
	}

	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return Proxy{}, fmt.Errorf("invalid port number: %v", err)
	}

	return Proxy{
		Protocol: protocol,
		Host:     hostParts[0],
		Port:     port,
		Username: username,
		Password: password,
	}, nil
}

func (b *Browser) Work(ctx context.Context, taskChan <-chan func()) {
	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return // Channel closed, exit gracefully
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Log.Sugar().Errorf("Panic recovered in browser task: %v", r)
					}
				}()

				task()
			}()

		case <-ctx.Done():
			return
		}
	}
}
