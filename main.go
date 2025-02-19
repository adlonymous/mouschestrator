package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"orchestrator/manager"
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

	workers := []string{fmt.Sprintf("%s:%d", host, port)}
	m := manager.New(workers)

	for i := 0; i < 3; i++ {
		t := task.Task{
			ID:    uuid.New(),
			Name:  fmt.Sprintf("test-container-%d", i),
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID:    uuid.New(),
			State: task.Running,
			Task:  t,
		}
		m.AddTask(te)
		m.SendWork()
	}

	go func() {
		for {
			fmt.Printf("[Manager] Updating tasks from %d workers\n", len(m.Workers))
			m.UpdateTasks()
			time.Sleep(15 * time.Second)
		}
	}

	for {
		for _, t := range m.TaskDb {
			fmt.Printf("[Manager] Task: id: %d, state: %d\n", t.ID, t.State)
			time.Sleep(15 * time.Second)
		}
	}
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
