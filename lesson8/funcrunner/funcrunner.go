package funcrunner

import (
	"fmt"
)

// waiting for all goroutines finish
func waitForCompleteTasks(taskCount int, taskChannel chan int) {
	for i := 0; i < taskCount; i++ {
		if _, ok := <-taskChannel; !ok {
			return
		}
	}
	close(taskChannel)
}

// Run concurrent tasks
func Run(tasks []func() error, N int, M int) error {
	// running task pool
	currentTaskChannel := make(chan int, N)
	// all task pool (running + finished)
	allTaskChannel := make(chan int, len(tasks))
	// error pool
	errorChannel := make(chan error, len(tasks))
	var currentTask int

	defer func() {
		waitForCompleteTasks(currentTask, allTaskChannel)
		close(currentTaskChannel)
		close(errorChannel)
	}()

	for currentTask < len(tasks) {

		if len(errorChannel) >= M {
			return fmt.Errorf("Max number of errors reached")
		}

		currentTaskChannel <- 1

		go func(currentTask int) {
			allTaskChannel <- 1
			err := tasks[currentTask]()
			if err != nil {
				errorChannel <- err
			}
			// decrease for clear N-task pool
			<-currentTaskChannel
		}(currentTask)

		currentTask++
	}

	waitForCompleteTasks(currentTask, allTaskChannel)
	if len(errorChannel) >= M {
		return fmt.Errorf("Max number of errors reached")
	}

	return nil
}
