package timer

import (
	"context"
	"time"
)

// StartTimer 定时任务
func StartTimer(ctx context.Context, fn func(t time.Time), hour int, min int) {
	startTimer(ctx, fn, hour, min, 0)
}

func startTimer(ctx context.Context, fn func(t time.Time), hour int, min int, sec int) {
	t := time.NewTimer(60 * time.Second)
	defer t.Stop()
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
		strat := next.Sub(now)

		if strat < 0 {
			next = next.AddDate(0, 0, 1)
			strat = next.Sub(now)
		}

		t.Reset(strat)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			fn(next)
		}
	}
}
