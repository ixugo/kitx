package task

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	// 创建一个并发量最大时为 10 的任务队列
	q := NewQueue(10, 1)

	var finishCount, errCount int32

	for i := 0; i < 100; i++ {
		err := q.Start(strconv.Itoa(i), func() {
			atomic.AddInt32(&finishCount, 1)
			time.Sleep(100 * time.Millisecond)
		})
		if err != nil {
			atomic.AddInt32(&errCount, 1)
			t.Log(err)
		}
	}

	q.Close()
	t.Logf("完成任务数量 %d\t丢弃任务数量 %d\n", finishCount, errCount)

}

func TestClose(t *testing.T) {
	q := NewQueue(10, 1)

	_ = q.Start("1", func() {
		fmt.Println("test")
	})

	q.Close()

}
