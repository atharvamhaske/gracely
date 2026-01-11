package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func slowJob1(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting job1 %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("Finished job1 %s\n", name)
}

func slowJob2(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting job2 %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("Finished job2 %s\n", name)
}

func slowJob3(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting job3 %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("Finished job3 %s\n", name)
}

func consumer(ctx context.Context, jobQueue chan string, doneChan chan interface{}) {
	wg := &sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			fmt.Println("writing to done channel")
			doneChan <- struct{}{}
			log.Println("Done, shutting down the consumer")
			return
		case job := <-jobQueue:
			wg.Add(3)
			go slowJob1(job, wg)
			go slowJob2(job, wg)
			go slowJob3(job, wg)
		}
	}
}
