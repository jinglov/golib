package logger

import (
	"sync"
	"testing"
)

func init() {
	SetLogFile("./log", "test", 5, M)
}

func TestSetLogFile(t *testing.T) {
	SetLogFile("./log", "test", 1, M)
}

func TestDebug(t *testing.T) {
	var sw sync.WaitGroup
	for i := 0; i < 100000; i++ {
		sw.Add(1)
		go func(i int) {
			defer sw.Done()
			Debug("this is debug:", i)
			Error("this is error:", i)
		}(i)
	}
	sw.Wait()
}
