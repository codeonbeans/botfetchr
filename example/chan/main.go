package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Regular Counter (unsafe for concurrent use)
type UnsafeCounter struct {
	Value int
}

// Atomic Counter (safe for concurrent use)
type AtomicCounter struct {
	Value int64
}

// Mutex-protected Counter (also safe)
type MutexCounter struct {
	Value int
	mu    sync.Mutex
}

func main() {
	fmt.Println("=== Demonstrating Race Conditions ===")

	// 1. Unsafe concurrent access (will have race conditions)
	fmt.Println("\n1. Unsafe Counter:")
	unsafeCounter := &UnsafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				unsafeCounter.Value++ // RACE CONDITION!
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Expected: 100000, Got: %d (likely wrong due to race condition)\n", unsafeCounter.Value)

	// 2. Using atomic operations (safe)
	fmt.Println("\n2. Atomic Counter:")
	atomicCounter := &AtomicCounter{}
	wg = sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&atomicCounter.Value, 1) // SAFE!
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Expected: 100000, Got: %d (correct with atomics)\n", atomicCounter.Value)

	// 3. Using mutex (also safe but slower)
	fmt.Println("\n3. Mutex Counter:")
	mutexCounter := &MutexCounter{}
	wg = sync.WaitGroup{}

	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutexCounter.mu.Lock()
				mutexCounter.Value++ // SAFE with mutex
				mutexCounter.mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Expected: 100000, Got: %d (correct with mutex)\n", mutexCounter.Value)
	fmt.Printf("Time taken: %v\n", time.Since(start))
}
