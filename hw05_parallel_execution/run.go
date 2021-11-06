package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	workers := make(chan int, n)
	var wg sync.WaitGroup
	var errorTasksCount int32

	for i, t := range tasks {
		if int(errorTasksCount) >= m && m > 0 {
			close(workers)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}

		workers <- i
		// fmt.Println("worker used: ", i)

		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			err := t()
			if err != nil {
				atomic.AddInt32(&errorTasksCount, 1)
			}
			// fmt.Println("worker free: ", <-workers)
			<-workers
		}(t)
	}

	wg.Wait()
	return nil
}
