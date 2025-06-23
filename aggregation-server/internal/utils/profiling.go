package utils

import (
	"log"
	"runtime"
	"time"
)

func PrintMemoryStats(prefix string) {
	if prefix == "" {
		prefix = "Memory - "
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("%s - Alloc: %v, TotalAlloc: %v, Sys: %v\n", prefix, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys))
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintTimeForFunction(name string, f func()) {
	if name == "" {
		name = "Function"
	}
	start := time.Now()
	f()
	duration := time.Since(start)
	log.Printf("%s-time: %v\n", name, duration)
}

func PrintTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s-time: %s", name, elapsed)
}
