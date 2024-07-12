package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"backend/auth"
	"backend/database"
	"backend/middleware"
	"backend/task"
	"backend/work"
)

func main() {

	r := mux.NewRouter()

	// Inisialisasi database
	if err := database.InitializeDB("root:@tcp(127.0.0.1:3306)/todo_golang"); err != nil {
		panic(err)
	}
	defer database.DB.Close()

	// Konfigurasi middleware CORS
	var corsOpts = cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Router login
	r.HandleFunc("/api/token", auth.HandleToken).Methods("POST")
	r.HandleFunc("/logout", auth.HandleLogout).Methods("POST")

	// Router task (dengan middleware otentikasi)
	r.HandleFunc("/api/task", task.CreateTask).Methods("POST").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(task.CreateTask))))
	r.HandleFunc("/api/task", task.GetTasks).Methods("GET").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(task.GetTasks))))
	r.HandleFunc("/api/task/{id}", task.GetTaskByID).Methods("GET").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(task.GetTaskByID))))
	r.HandleFunc("/api/task/{id}", task.UpdateTask).Methods("PUT").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(task.UpdateTask))))
	r.HandleFunc("/api/task/{id}", task.DeleteTask).Methods("DELETE").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(task.DeleteTask))))

	r.HandleFunc("/api/work/{id}", work.UpdateWorkLog).Methods("PUT").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(work.UpdateWorkLog))))
	r.HandleFunc("/api/work", work.GetWorkLogs).Methods("GET").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(work.GetWorkLogs))))
	r.HandleFunc("/api/work/today", work.GetTodayWorkLogs).Methods("GET").Handler(middleware.AuthMiddleware(corsOpts.Handler(http.HandlerFunc(work.GetTodayWorkLogs))))

    // go func() {
    //     for {
    //         // Calculate time until next midnight
    //         now := time.Now()
    //         nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    //         if nextMidnight.Before(now) {
    //             nextMidnight = nextMidnight.AddDate(0, 0, 1) 
    //         }
    //         durationUntilMidnight := nextMidnight.Sub(now)

    //         // Sleep until midnight
    //         time.Sleep(durationUntilMidnight)

    //         // Refresh tasks after midnight
    //         work.AssignDailyTasks(nil, nil)
    //     }
    // }()

	go func() {
		for {
			// Refresh tasks after midnight with dummy w and r
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			work.AssignDailyTasks(w, r) 
			time.Sleep(100 * time.Minute) // Tunggu 1 menit (or another short interval) for testing
			log.Printf("Task assigned and refreshed")
		}
	}()

	// Serve server dengan middleware CORS
	port := "8080"

	// Convert the port string to an integer
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Error converting port to integer:", err)
	}

	// Start the server with the converted port
	fmt.Println("Server running on port:", port)
	http.ListenAndServe(":"+strconv.Itoa(portInt), corsOpts.Handler(r))
}
