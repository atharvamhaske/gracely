package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type CustomHandler struct {
	wg *sync.WaitGroup
}

func NewCustomHandler(wg *sync.WaitGroup) *CustomHandler {
	return &CustomHandler{
		wg: wg,
	}
}

func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	jobName := vars["jobName"]

	fmt.Fprintf(w, "job %s started \n", jobName)

	h.wg.Add(3)
	go job1(jobName, h.wg)
	go job2(jobName, h.wg)
	go job3(jobName, h.wg)
}

func job1(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("starting the job1 for %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("finished job1 for %s\n", name)
}

func job2(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("starting the job2 for %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("finished job2 for %s\n", name)
}

func job3(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("starting the job3 for %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("finished job3 for %s\n", name)
}

func main() {
	wg := &sync.WaitGroup{}

	customHandler := NewCustomHandler(wg)

	r := mux.NewRouter()

	r.Handle("/{jobName}", customHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	tChan := make(chan os.Signal, 1)

	signal.Notify(tChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-tChan //blocks here until interrupted
		log.Printf("SIGTERM received, Shutdown process initiated\n")
		httpServer.Shutdown(context.Background())
	}()

	err := httpServer.ListenAndServe()
	if err != nil {
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP Server Closed with error: %v\n", err)
		}
		log.Printf("HTTP server shut down")
	}

	log.Println("Waiting for all jobs to finish")
	wg.Wait()
	log.Println("All Jobs finished")
}
