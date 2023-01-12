package timer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestStartTimer(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	startTimer(ctx, func(d time.Time) {
	}, now.Hour(), now.Minute(), now.Second()+1)
}

func BenchmarkStartTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		b.StartTimer()
		startTimer(ctx, func(ti time.Time) {}, 0, 0, 0)
		b.StopTimer()
		cancel()
	}
}
