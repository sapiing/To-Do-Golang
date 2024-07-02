package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"

	"backend/auth"
	"backend/task"
	"backend/work"
	"backend/database"
	"backend/middleware"
)



func main() {

	// inisiasi db
	if err := database.InitializeDB("root:@tcp(127.0.0.1:3306)/todo_golang"); err != nil {
        panic(err)
    }
    defer database.DB.Close()

    r := mux.NewRouter() 

	// router login
    r.HandleFunc("/api/token", auth.HandleToken).Methods("POST") 	

	// router task
    r.HandleFunc("/api/task", task.CreateTask).Methods("POST").Handler(middleware.AuthMiddleware(http.HandlerFunc(task.CreateTask)))
    r.HandleFunc("/api/task", task.GetTasks).Methods("GET").Handler(middleware.AuthMiddleware(http.HandlerFunc(task.GetTasks)))
    r.HandleFunc("/api/task/{id}", task.GetTaskByID).Methods("GET").Handler(middleware.AuthMiddleware(http.HandlerFunc(task.GetTaskByID)))
    r.HandleFunc("/api/task/{id}", task.UpdateTask).Methods("PUT").Handler(middleware.AuthMiddleware(http.HandlerFunc(task.UpdateTask)))
    r.HandleFunc("/api/task/{id}", task.DeleteTask).Methods("DELETE").Handler(middleware.AuthMiddleware(http.HandlerFunc(task.DeleteTask)))

	// router work
    r.HandleFunc("/api/work", work.CreateWork).Methods("POST").Handler(middleware.AuthMiddleware(http.HandlerFunc(work.CreateWork)))
    r.HandleFunc("/api/work", work.GetWorks).Methods("GET").Handler(middleware.AuthMiddleware(http.HandlerFunc(work.GetWorks)))
    r.HandleFunc("/api/work/{id}", work.GetWorkByID).Methods("GET").Handler(middleware.AuthMiddleware(http.HandlerFunc(work.GetWorkByID)))
    r.HandleFunc("/api/work/{id}", work.UpdateWork).Methods("PUT").Handler(middleware.AuthMiddleware(http.HandlerFunc(work.UpdateWork)))
    r.HandleFunc("/api/work/{id}", work.DeleteWork).Methods("DELETE").Handler(middleware.AuthMiddleware(http.HandlerFunc(work.DeleteWork)))



	// serve server
    port := ":8080"
    fmt.Println("Server berjalan di port", port)
    http.ListenAndServe(port, r)
}