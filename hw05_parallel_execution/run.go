package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	workers := make(chan int, n)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorTasksCount int

	for i, t := range tasks {
		mu.Lock()
		errorTasksCountCheck := errorTasksCount
		mu.Unlock()

		if errorTasksCountCheck >= m && m > 0 {
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
				mu.Lock()
				errorTasksCount++
				mu.Unlock()
			}
			// fmt.Println("worker free: ", <-workers)
			<-workers
		}(t)
	}

	wg.Wait()
	return nil
}
