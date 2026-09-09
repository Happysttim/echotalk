package timer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TimeTicker struct {
	ticker *time.Ticker
	ctx    context.Context
	stop   context.CancelFunc
}

func NewTimeTicker(tick time.Duration) *TimeTicker {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	return &TimeTicker{
		ticker: time.NewTicker(tick),
		ctx:    ctx,
		stop:   stop,
	}
}

func (ticker *TimeTicker) Start(handler func(context.Context)) {
	go ticker.start(handler)
}

func (ticker *TimeTicker) start(handler func(context.Context)) {
	for {
		select {
		case <-ticker.ticker.C:
			handler(ticker.ctx)
		case <-ticker.ctx.Done():
			ticker.ticker.Stop()
			return
		}
	}
}

func (ticker *TimeTicker) Stop() {
	ticker.stop()
}
