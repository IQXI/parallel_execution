package main

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

func executor(functions []func() error, N, errCount int) {
	if len(functions) < N {
		N = len(functions)
	}

	channel := make(chan struct{}, N)
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	wg.Add(len(functions))

	for i := 0; i < len(functions); i++ {
		log.Printf("error count: %v", errCount)
		if errCount <= 0 {
			log.Printf("\nError count is expired")
			os.Exit(1)
		}

		channel <- struct{}{}

		go func(f func() error, errC *int) {
			err := f()
			if err != nil {
				log.Print(err)
				mutex.Lock()
				*errC--
				mutex.Unlock()
			}
			<-channel
			wg.Done()

		}(functions[i], &errCount)

	}
	wg.Wait()

}

func one() error {
	log.Printf("One")
	time.Sleep(time.Millisecond * 450)
	return nil
}

func two() error {
	log.Printf("Two")
	time.Sleep(time.Millisecond * 500)
	return nil
}

func three() error {
	log.Printf("Three")
	time.Sleep(time.Millisecond * 550)
	return errors.New("Some error!")
}

func main() {
	functions := []func() error{one, two, three, one, two, three, two, three, one, two, three, two, three, one, two, three}
	executor(functions, 3, 2)
}
