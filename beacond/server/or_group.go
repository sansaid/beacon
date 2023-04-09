package server

import (
	"context"
)

type OrGroup struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Error  error
}

type OrGroupManager interface {
	Go(func() error)
	Wait() error
}

func NewOrGroup() OrGroupManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &OrGroup{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

// Error if any one of the goroutines in an or group error
func (o *OrGroup) Go(routine func() error) {
	go func() {
		defer o.Cancel()

		err := routine()

		if err != nil {
			o.Error = err
		}
	}()
}

func (o *OrGroup) Wait() error {
	<-o.Ctx.Done()
	return o.Error
}
