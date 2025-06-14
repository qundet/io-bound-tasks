package main

import (
	"fmt"
	h "io-bound-tasks/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			http.NotFound(w, r)
			return
		}
		h.CreateTask(w, r)
	})

	http.HandleFunc("/tasks/", h.TaskHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
