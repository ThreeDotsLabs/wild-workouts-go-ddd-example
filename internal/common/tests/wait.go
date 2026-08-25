package tests

import (
	"context"
	"net"
	"time"
)

func WaitForPort(address string) bool {
	waitChan := make(chan struct{})

	go func() {
		dialer := &net.Dialer{Timeout: time.Second}
		for {
			conn, err := dialer.DialContext(context.Background(), "tcp", address)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if conn != nil {
				_ = conn.Close()
				waitChan <- struct{}{}
				return
			}
		}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-waitChan:
		return true
	case <-timeout:
		return false
	}
}
