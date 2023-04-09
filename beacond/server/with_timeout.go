package server

import (
	"fmt"
	"time"
)

func withTimeout(f func() error, timeout time.Duration, errorMsg string) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Reset(timeout)

	var finished chan struct{}

	var err error
	err = nil

	go func() {
		defer func() { finished <- struct{}{} }()

		err = f()
	}()

	for {
		select {
		case <-ticker.C:
			return fmt.Errorf(errorMsg)
		case <-finished:
			return err
		}
	}
}
