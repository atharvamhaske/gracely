package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type CustomHandle struct {
	jobQueue chan string
}

func CustomHandler(jobQueue chan string) *CustomHandle {
	return &CustomHandle{
		jobQueue: jobQueue,
	}
}

func (h *CustomHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["jobName"]

	h.jobQueue <- jobName

	fmt.Fprintf(w, "job %s started", jobName)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	jobQueue := make(chan string)

	customHandler := CustomHandler(jobQueue) //passing newly created jobQ channel in CustomHandler

	r := mux.NewRouter()
	r.Handle("/{jobName}", customHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	doneChan := make(chan any)
	go consumer(ctx, jobQueue, doneChan)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Printf("HTTP server closed with: %v\n", err)
			}
			log.Printf("HTTP server shut down")
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	<-termChan
	log.Println("SIGTERM received. Shutdown process initiated")

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	cancel()

	log.Println("waiting consumer to finish its jobs")
	<-doneChan
	log.Println("done. returning.")
}
