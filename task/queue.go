// Package task 创建一个 goroutine，线性执行任务
package task

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Queue ...
type Queue struct {
	wg sync.WaitGroup
	ch chan func()
}

// NewQueue 填入 channel 容量和 goroutine 数量
// 线性执行，请设置 worker=1
func NewQueue(cap, worker int) *Queue {
	q := Queue{
		ch: make(chan func(), cap),
	}

	// 设置上下限
	if worker <= 0 {
		worker = 1
	}
	if worker > runtime.NumCPU() {
		worker = runtime.NumCPU()
	}
	q.wg.Add(worker)
	for i := 0; i < worker; i++ {
		go func() {
			defer q.wg.Done()
			for fn := range q.ch {
				fn()
			}
		}()
	}
	return &q
}

// Start 往队列加一个新的任务
func (q *Queue) Start(sign string, fn func()) error {
	select {
	case q.ch <- fn:
	case <-time.After(100 * time.Millisecond):
		// default:
		return fmt.Errorf("队列已满，任务「%s」已丢弃", sign)
	}
	return nil
}

// Close 调用后，请勿再添加任务或重复关闭
func (q *Queue) Close() {
	close(q.ch)
	q.wg.Wait()
}
