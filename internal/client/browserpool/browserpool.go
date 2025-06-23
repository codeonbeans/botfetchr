package browserpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Client interface {
	UseBrowser(fn func(browser *Browser) error) error
	Close() error
}

type clientImpl struct {
	browsers   []*Browser
	taskChans  []chan func(browser *Browser) error
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	closed     int32
	roundRobin int64
}

type Config struct {
	Proxies       []string
	PoolSize      int
	TaskQueueSize int // Buffer size for task channels
}

func NewClient(cfg Config) (Client, error) {
	if cfg.PoolSize <= 0 {
		return nil, fmt.Errorf("pool size must be greater than 0")
	}
	if cfg.TaskQueueSize <= 0 {
		return nil, fmt.Errorf("task queue size must be greater than 0")
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &clientImpl{
		browsers:  make([]*Browser, 0, cfg.PoolSize),
		taskChans: make([]chan func(browser *Browser) error, 0, cfg.PoolSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Create browsers and their corresponding task channels
	for i := 0; i < cfg.PoolSize; i++ {
		var proxy string
		if i < len(cfg.Proxies) {
			proxy = cfg.Proxies[i]
		}

		browser, err := NewBrowser(proxy)
		if err != nil {
			// Clean up any already created browsers
			client.Close()
			return nil, fmt.Errorf("failed to create browser %d: %w", i, err)
		}

		taskChan := make(chan func(browser *Browser) error, cfg.TaskQueueSize)

		client.browsers = append(client.browsers, browser)
		client.taskChans = append(client.taskChans, taskChan)

		// Start worker goroutine for this browser
		client.wg.Add(1)
		go func(browser *Browser, taskChan chan func(browser *Browser) error) {
			defer client.wg.Done()
			browser.Work(ctx, taskChan)
		}(browser, taskChan)
	}

	return client, nil
}

func (c *clientImpl) UseBrowser(fn func(browser *Browser) error) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return fmt.Errorf("client is closed")
	}

	index := atomic.AddInt64(&c.roundRobin, 1) % int64(len(c.taskChans))

	// Create a channel to wait for task completion
	resultChan := make(chan error, 1)

	task := func(browser *Browser) error {
		err := fn(browser)
		resultChan <- err
		return err
	}

	select {
	case c.taskChans[index] <- task:
		// Wait for task completion
		err := <-resultChan
		return err
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		return fmt.Errorf("task queue is full")
	}
}

func (c *clientImpl) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil // Already closed
	}

	// Cancel context to signal all workers to stop
	c.cancel()

	// Close all task channels
	for _, taskChan := range c.taskChans {
		close(taskChan)
	}

	// Wait for all workers to finish
	c.wg.Wait()

	// Close all browsers
	var errs []error
	for _, browser := range c.browsers {
		if err := browser.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing browsers: %v", errs)
	}

	return nil
}
