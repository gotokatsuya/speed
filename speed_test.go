package speed

import (
	"sync"
	"testing"
	"time"
)

func TestParallelSpeedLog(t *testing.T) {
	EnableLogger()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		speedLogger := NewLogger("test-parallel-log-1").Begin()
		defer speedLogger.End()
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		speedLogger := NewLogger("test-parallel-log-2").Begin()
		defer speedLogger.End()
		time.Sleep(2 * time.Second)
	}()

	wg.Wait()
}

var globalSpeedLogger *Logger

func TestGlobalSpeedLog(t *testing.T) {
	EnableLogger()

	globalSpeedLogger = NewLogger("test-global-log")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		speedLogger := globalSpeedLogger.Copy().Begin()
		defer speedLogger.End()
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		speedLogger := globalSpeedLogger.Copy().Begin()
		defer speedLogger.End()
		time.Sleep(2 * time.Second)
	}()

	wg.Wait()
}
