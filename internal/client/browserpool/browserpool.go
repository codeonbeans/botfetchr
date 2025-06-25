package browserpool

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const resetThreshold = 1000000

type Client interface {
	UseBrowser(fn func(ctx context.Context, browser *Browser) error) error
	Close() error
}

type clientImpl struct {
	ctx    context.Context
	wg     sync.WaitGroup
	cancel context.CancelFunc

	taskTimeout time.Duration // Timeout for tasks
	browsers    []*Browser
	taskChans   []chan func()
	closed      int32
	roundRobin  uint64 // For round-robin browser selection
}

type Config struct {
	Headless      bool          // Whether to run browsers in headless mode
	Proxies       []string      // List of proxies to use for browsers, can be empty
	PoolSize      int           // Number of browsers in the pool
	TaskQueueSize int           // Buffer size for task channels
	TaskTimeout   time.Duration // Timeout for tasks
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
		ctx:         ctx,
		cancel:      cancel,
		taskTimeout: cfg.TaskTimeout,
		browsers:    make([]*Browser, 0, cfg.PoolSize),
		taskChans:   make([]chan func(), 0, cfg.PoolSize),
	}

	go func() {
		// Wait for interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("Received shutdown signal, closing all browsers...")
		if err := client.Close(); err != nil {
			fmt.Printf("Error closing browsers: %v\n", err)
		}
		os.Exit(0)
	}()

	// Create browsers and their corresponding task channels
	for i := 0; i < cfg.PoolSize; i++ {
		var proxy string
		if i < len(cfg.Proxies) {
			proxy = cfg.Proxies[i]
		}

		browser, err := NewBrowser(cfg.Headless, proxy)
		if err != nil {
			// Clean up any already created browsers
			client.Close()
			return nil, fmt.Errorf("failed to create browser %d: %w", i, err)
		}

		taskChan := make(chan func(), cfg.TaskQueueSize)

		client.browsers = append(client.browsers, browser)
		client.taskChans = append(client.taskChans, taskChan)

		// Start worker goroutine for this browser
		client.wg.Add(1)
		go func() {
			defer client.wg.Done()
			browser.Work(ctx, taskChan)
		}()
	}

	return client, nil
}

func (c *clientImpl) UseBrowser(fn func(ctx context.Context, browser *Browser) error) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return fmt.Errorf("client is closed")
	}

	index := atomic.AddUint64(&c.roundRobin, 1) % uint64(len(c.taskChans))
	if atomic.LoadUint64(&c.roundRobin) > resetThreshold {
		atomic.StoreUint64(&c.roundRobin, 0)
	}

	browser := c.browsers[index]
	taskChan := c.taskChans[index]

	// Create a channel to wait for task completion
	resultChan := make(chan error, 3)

	taskChan <- func() {
		// Make sure to recover from any panic in the task
		// avoid stuck result channel
		defer func() {
			if r := recover(); r != nil {
				resultChan <- fmt.Errorf("panic recovered: %v", r)
			}
			close(resultChan)
		}()

		// Timeout here not did not solve something, the task still be stucked btw. Should handle timeout in the task itself tho
		// go func() {
		// 	time.Sleep(c.taskTimeout)
		// 	resultChan <- fmt.Errorf("task timed out")
		// }()

		resultChan <- fn(c.ctx, browser)
	}

	select {
	case err := <-resultChan:
		return err
	case <-c.ctx.Done():
		return c.ctx.Err()
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
