package funcrunner

import (
	"fmt"
)

// Run concurrent tasks
func Run(tasks []func() error, N int, M int) error {
	// start task pool
	taskChannel := make(chan int, len(tasks))
	// error pool
	errorChannel := make(chan error, len(tasks))

	// routine count
	routineNum := N
	if len(tasks) < N {
		routineNum = len(tasks)
	}

	// quit task pool
	quitChannel := make(chan int, routineNum)
	for i := 0; i < routineNum; i++ {
		quitChannel <- i
	}

	defer func() {
		close(errorChannel)
		close(quitChannel)
	}()

	// feel pool with tasks id-s
	for i := 0; i < len(tasks); i++ {
		taskChannel <- i
	}
	close(taskChannel)

	for i := 0; i < routineNum; i++ {

		go func() {

			defer func() {
				<-quitChannel
			}()

			for {
				if len(errorChannel) >= M {
					return
				}
				taskID, ok := <-taskChannel
				if !ok {
					return
				}
				err := tasks[taskID]()
				if err != nil {
					errorChannel <- err
				}
			}

		}()

	}

	for len(quitChannel) > 0 {
	}

	if len(errorChannel) >= M {
		return fmt.Errorf("Max number of errors reached")
	}

	return nil
}
