package timer

import (
	"context"
	"time"
)

type TimerState int

const (
	Stopped TimerState = iota
	Running
)

type BasicTimer struct {
	State       TimerState
	StartTime   time.Time
	LastElapsed time.Duration

	cancelFunc context.CancelFunc
	ctx        context.Context
}

func NewBasicTimer(lastTime time.Duration) *BasicTimer {
	return &BasicTimer{
		LastElapsed: lastTime,
		State:       Stopped,
	}
}

func (bt *BasicTimer) Start(ch chan time.Duration, notifyPoints ...time.Duration) {
	if ch != nil && len(notifyPoints) > 0 {
		bt.ctx, bt.cancelFunc = context.WithCancel(context.Background())
		go NotifyAtPoints(bt.ctx, ch, notifyPoints...)
	}
	bt.StartTime = time.Now()
	bt.State = Running
}

func (bt *BasicTimer) Peek() time.Duration {
	if bt.State != Running {
		return bt.LastElapsed
	}
	return time.Since(bt.StartTime)
}

func (bt *BasicTimer) Stop() time.Duration {
	if bt.cancelFunc != nil {
		bt.cancelFunc()
		bt.cancelFunc = nil
	}
	bt.LastElapsed = time.Since(bt.StartTime)
	bt.State = Stopped
	return bt.LastElapsed
}

func NotifyAtPoints(ctx context.Context, ch chan time.Duration, notifyPoints ...time.Duration) {
	for _, point := range notifyPoints {
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		case <-time.After(point - time.Since(start)):
			ch <- point
		}
	}
}
