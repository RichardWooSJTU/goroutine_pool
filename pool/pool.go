package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type GoroutinePool struct {
	cap          uint64
	runningTasks uint64

	chTasks chan *Task
	status  uint

	sync.Mutex
}

func (pool *GoroutinePool) GetCap() uint64 {
	return pool.cap
}

func (pool *GoroutinePool) GetRunningTasks() uint64 {
	return pool.runningTasks
}

func (pool *GoroutinePool) GetStatus() uint {
	return pool.status
}

func (pool *GoroutinePool) SetStatus(stat uint) {
	pool.Lock()
	defer pool.Unlock()
	pool.status = stat
}
var ErrInvalidCap = errors.New("invalid pool cap")

const (
	RUNNING = 1
	STOP    = 0
)

func New(cap uint64) (*GoroutinePool, error) {
	if cap <= 0 {
		return nil, ErrInvalidCap
	}
	return &GoroutinePool{
		cap:          cap,
		runningTasks: 0,
		status:       RUNNING,
		chTasks:      make(chan *Task, cap),
	}, nil
}

func (pool *GoroutinePool) Put(task *Task) {
	pool.Lock()
	defer pool.Unlock()
	if pool.status == STOP {
		return
	}
	if pool.runningTasks < pool.cap {
		pool.run()
	}
	pool.chTasks <- task
}

func (pool *GoroutinePool) Close() {
	pool.SetStatus(STOP)
	for len(pool.chTasks) > 0 {
		time.Sleep(1e6)
	}
	close(pool.chTasks)
}

func (pool *GoroutinePool) decrement(i *uint64) {
	atomic.AddUint64(i, ^uint64(0))//0xffffffff 加上一个uint就是-1 因为这个数字比0小1
}

func (pool *GoroutinePool) run() {
	pool.runningTasks++

	go func() {
		defer func() {
			pool.decrement(&pool.runningTasks)
		}()
		for {
			select {
			case task, ok := <-pool.chTasks:
				if !ok {
					return //如果队列已经被关闭了，go程退出，如果队列为空select会把它阻塞
				}
				task.Handler(task.Parameters...)
			}
		}
	}()
}
