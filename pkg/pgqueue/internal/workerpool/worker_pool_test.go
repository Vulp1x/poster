package workerpool

import (
	"testing"
	"time"
)

func TestWorkerPoolResize(t *testing.T) {
	t.Parallel()
	pool := New()
	size := int32(8)

	check := func(size int32) {
		if s := pool.Size(); s != size {
			t.Fatalf("Size: expected [%v], got [%v]", size, s)
		}
		if s := pool.RestingCount(); s != size {
			t.Fatalf("RestingCount: expected [%v], got [%v]", size, s)
		}
	}

	pool.Resize(size)
	check(size)

	pool.Resize(2 * size)
	check(2 * size)

	pool.Resize(size)
	check(size)

	pool.CloseAndWait()
}

func TestWorkerPoolTasks(t *testing.T) {
	t.Parallel()
	pool := New()

	size := int32(4)
	pool.Resize(size)

	chs, tasks := createTasks(size)
	for i := int32(0); i < size; i++ {
		pool.Push(tasks[i])
		time.Sleep(50 * time.Millisecond)
		if s := pool.RestingCount(); s != size-(i+1) {
			t.Fatalf("RestingCount: expected [%v], got [%v]", size-(i+1), s)
		}
	}
	for i := int32(0); i < size; i++ {
		close(chs[i])
		time.Sleep(50 * time.Millisecond)
		if s := pool.RestingCount(); s != i+1 {
			t.Fatalf("RestingCount: expected [%v], got [%v]", i+1, s)
		}
	}

	pool.CloseAndWait()
}

func newTask() (chan struct{}, Task) {
	ch := make(chan struct{})
	f := func() {
		<-ch
	}
	return ch, f
}

func createTasks(n int32) ([]chan struct{}, []Task) {
	chs := make([]chan struct{}, n)
	tasks := make([]Task, n)
	for i := int32(0); i < n; i++ {
		chs[i], tasks[i] = newTask()
	}
	return chs, tasks
}
