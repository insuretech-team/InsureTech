package brokerutil

import (
	"context"
	"net"
	"time"
)

const defaultProbeTimeout = 750 * time.Millisecond

// ReachableBrokers returns the subset of brokers that accept a TCP connection
// within the timeout. If none are reachable, callers should fall back to the
// original list so the underlying Kafka client can return the real error.
func ReachableBrokers(brokers []string, timeout time.Duration) (reachable []string, unreachable []string) {
	if len(brokers) == 0 {
		return nil, nil
	}
	if timeout <= 0 {
		timeout = defaultProbeTimeout
	}

	type result struct {
		index     int
		reachable bool
	}

	results := make(chan result, len(brokers))
	for i, broker := range brokers {
		go func(idx int, addr string) {
			results <- result{index: idx, reachable: canDial(addr, timeout)}
		}(i, broker)
	}

	status := make([]bool, len(brokers))
	for range brokers {
		res := <-results
		status[res.index] = res.reachable
	}

	for i, broker := range brokers {
		if status[i] {
			reachable = append(reachable, broker)
		} else {
			unreachable = append(unreachable, broker)
		}
	}

	return reachable, unreachable
}

func canDial(address string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
