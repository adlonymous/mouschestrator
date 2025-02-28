package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"orchestrator/manager"
	"orchestrator/task"
	"orchestrator/worker"
	"os"
	"strconv"
)

func main() {
	whost := os.Getenv("MOUSCHESTRATOR_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("MOUSCHESTRATOR_WORKER_PORT"))

	mhost := os.Getenv("MOUSCHESTRATOR_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("MOUSCHESTRATOR_MANAGER_PORT"))

	fmt.Println("Starting Cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	fmt.Println("Starting Cube manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()

	mapi.Start()
}
