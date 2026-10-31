package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-ping/ping"
)

const healthCheckInterval = 5 * time.Second

// HEALTH CHECK
func runHealthCheck(targets []Target) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		results := make([]CheckResult, len(targets))

		var wg sync.WaitGroup

		for i, target := range targets {
			wg.Add(1)

			go func(index int, t Target) {
				defer wg.Done()

				results[index] = checkTarget(t)
			}(i, target)
		}

		wg.Wait()

		return healthCheckMsg{
			results:  results,
			duration: time.Since(start),
		}
	}
}

func checkTarget(target Target) CheckResult {
	start := time.Now()

	result := CheckResult{
		Target: target,
	}

	pinger, err := ping.NewPinger(target.URL)
	if err != nil {
		result.Status = "DOWN"
		result.Error = err.Error()
		result.Latency = time.Since(start)
		return result
	}

	pinger.Count = 1
	pinger.Timeout = 4 * time.Second
	pinger.Interval = time.Second
	pinger.SetPrivileged(true)

	err = pinger.Run()
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = "DOWN"
		result.Error = err.Error()
		return result
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		result.Status = "UP"
	} else {
		result.Status = "DOWN"
		result.Error = fmt.Sprintf("packet loss: %.0f%%", stats.PacketLoss)
	}

	return result
}

func normalizeTarget(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}

	if strings.Contains(raw, "://") {
		raw = strings.TrimPrefix(raw, "http://")
		raw = strings.TrimPrefix(raw, "https://")
	}

	if strings.ContainsAny(raw, "/?#") {
		return "", false
	}

	return raw, true
}
