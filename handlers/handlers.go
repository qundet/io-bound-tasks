package handlers

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	"net/http"
	"encoding/json"
    "github.com/google/uuid"
	tm "io-bound-tasks/internal/task"

)

var (
	tasks      = make(map[uuid.UUID]*tm.Task)
	tasksMutex = &sync.Mutex{}
)

func simulateTask() {

	rand.Seed(time.Now().UnixNano())

    min := 3 * time.Minute
    max := 5 * time.Minute

    delta := time.Duration(rand.Int63n(int64(max - min))) + min

    fmt.Println("Do some smart stuff :)...")
    time.Sleep(delta)
    fmt.Println("Done")
}


func CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &tm.Task{
		ID:        uuid.New(),
		Status:    tm.StatusCreated,
		CreatedAt: time.Now(),
	}

	tasksMutex.Lock()
	tasks[task.ID] = task
	tasksMutex.Unlock()

	go func() {
		tasksMutex.Lock()
		task.Status = tm.StatusRunning
		tasksMutex.Unlock()

		defer func() {
			tasksMutex.Lock()
			task.CompletedAt = time.Now()

			tasksMutex.Unlock()
		}()

		simulateTask()
		task.Status = tm.StatusCompleted
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"id": task.ID})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateTask(w, r)
	case http.MethodGet:
		getTaskStatus(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Path[len("/tasks/"):]

    taskID, err := uuid.Parse(taskIDStr)
    if err != nil {
        http.Error(w, "Invalid UUID format", http.StatusBadRequest)
        return
    }
	tasksMutex.Lock()
	task, exists := tasks[taskID]
	tasksMutex.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task.ToJSON())
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Path[len("/tasks/"):]

    taskID, err := uuid.Parse(taskIDStr)
    if err != nil {
        http.Error(w, "Invalid UUID format", http.StatusBadRequest)
        return
    }
    
	tasksMutex.Lock()
	_, exists := tasks[taskID]
	if exists {
		delete(tasks, taskID)
	}
	tasksMutex.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
