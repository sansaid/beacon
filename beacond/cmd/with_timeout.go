package cmd

import (
	"fmt"
	"time"
)

func withTimeout(f func() error, timeout time.Duration, errorMsg string) error {
	var finished chan struct{}

	var err error
	err = nil

	go func() {
		defer func() { finished <- struct{}{} }()

		err = f()
	}()

	for {
		select {
		case <-time.Tick(timeout):
			return fmt.Errorf(errorMsg)
		case <-finished:
			return err
		}
	}
}
