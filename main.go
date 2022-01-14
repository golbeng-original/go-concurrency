package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/golbeng-original/go-concurrency/workerpool"
)

func useWorkerPool() {

	resultStream := make(chan interface{})
	result2Stream := make(chan interface{})

	taskPool1 := workerpool.NewPool(10)
	taskPool1.Run()

	taskPool2 := workerpool.NewPool(2)
	taskPool2.Run()

	taskPool3 := workerpool.NewPool(7)
	taskPool3.Run()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 7; i++ {
			taskPool1.AddTask(func(args ...interface{}) {
				time.Sleep(time.Millisecond * 500)
				fmt.Println("Run Task - ", args[0])

				resultStream <- args[0]
			}, i)
		}

		taskPool1.Done()

		taskPool1.Wait()
		close(resultStream)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for result := range resultStream {

			taskPool2.AddTask(func(args ...interface{}) {
				time.Sleep(time.Second * 2)
				fmt.Println("Run Task2 - ", args[0])

				result2Stream <- args[0]
			}, result)

		}

		taskPool2.Done()
		taskPool2.Wait()
		close(result2Stream)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for result := range result2Stream {

			taskPool3.AddTask(func(args ...interface{}) {
				time.Sleep(time.Second)
				fmt.Println("Run Task3 - ", args[0])
			}, result)

		}

		taskPool3.Done()
		taskPool3.Wait()
	}()

	wg.Wait()

}

func main() {
	useWorkerPool()
}
