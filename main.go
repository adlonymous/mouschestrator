package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"orchestrator/task"
	"orchestrator/worker"
	"os"
	"strconv"
	"time"
)

func main() {
	host := os.Getenv("MOUSCHESTRATOR_HOST")
	port, _ := strconv.Atoi(os.Getenv("MOUSCHESTRATOR_PORT"))

	fmt.Println("Starting Mouschestrator Worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error while running task: %v", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently \n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
