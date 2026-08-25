package client

import (
	"context"
	"net"
	"time"
)

func waitForPort(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
				// continue
			}

			conn, err := new(net.Dialer).DialContext(context.Background(), "tcp", addr)
			if err == nil {
				_ = conn.Close()
				close(portAvailable)
				return
			}

			time.Sleep(time.Millisecond * 200)
		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
