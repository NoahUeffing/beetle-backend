package main

import (
	"fmt"
	"log"
	"net/http"

	"beetle/internal/config"
	"beetle/internal/domain"
	"beetle/internal/handler"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	dbConfig := config.NewDBConfig()
	fmt.Println("Host")
	fmt.Println(dbConfig.Host)
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize router and handler
	r := mux.NewRouter()
	taskHandler := handler.NewTaskHandler(db)

	// Define routes
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// Task routes
	r.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Start server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to the Beetle API"}`))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
