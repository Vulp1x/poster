package workerpool

import (
	"sync"
	"sync/atomic"
)

type Task func()

type Pool struct {
	tasks chan Task
	kill  chan struct{}

	size      int32
	busyCount int32

	mu sync.Mutex
	wg sync.WaitGroup
}

func New() *Pool {
	return &Pool{tasks: make(chan Task, 16), kill: make(chan struct{})}
}

func (p *Pool) Push(task Task) {
	p.tasks <- task
}

func (p *Pool) Resize(size int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.size < size {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}

	for p.size > size && p.size > 0 {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) RestingCount() int32 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.size - atomic.LoadInt32(&p.busyCount)
}

func (p *Pool) CloseAndWait() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.kill)
	p.wg.Wait()
}

func (p *Pool) Size() int32 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.size
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.tasks:
			p.executeTask(task)
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) executeTask(task Task) {
	atomic.AddInt32(&p.busyCount, 1)
	defer atomic.AddInt32(&p.busyCount, -1)

	task()
}
