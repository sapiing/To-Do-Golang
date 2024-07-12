package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/database"
	"backend/error" // Import package error
	"backend/models"

	"github.com/gorilla/mux"
)

type Task = models.Task

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		error.HandleError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate title is not empty (optional)
	if task.Title == "" {
		error.HandleError(w, http.StatusBadRequest, "Title is required", nil)
		return
	}

	// Check for duplicate title using prepared statement
	stmt, err := database.DB.Prepare("SELECT id FROM tasks WHERE title = ?")
	if err != nil {
		error.HandleError(w, http.StatusInternalServerError, "Failed to prepare statement", err)
		return
	}
	defer stmt.Close()

	var existingTaskID int
	err = stmt.QueryRow(task.Title).Scan(&existingTaskID)

	if err == sql.ErrNoRows {
		// If no duplicate, set the Date field to the current time
		task.Date = time.Now()

		// Insert the task with the Date field
		result, err := database.DB.Exec(
			"INSERT INTO tasks (title, description, completed, date) VALUES (?, ?, ?, ?)",
			task.Title, task.Description, task.Completed, task.Date.Format(time.RFC3339),
		)
		if err != nil {
			error.HandleError(w, http.StatusInternalServerError, "Failed to create task", err)
			return
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			error.HandleError(w, http.StatusInternalServerError, "Failed to get last insert ID", err)
			return
		}
		task.ID = int(lastInsertID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
		return
	} else if err != nil {
		error.HandleError(w, http.StatusInternalServerError, "Failed to check for duplicate task", err)
		return
	}

	// Duplicate found
	error.HandleError(w, http.StatusConflict, "Task with this title already exists", nil)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent) // Optional: Use 204 for OPTIONS
		return
	}

	var tasks []Task
	rows, err := database.DB.Query("SELECT id, title, description, completed, date FROM tasks")
	if err != nil {
		error.HandleError(w, http.StatusInternalServerError, "Error querying tasks", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var dateBytes []byte
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &dateBytes); err != nil {
			error.HandleError(w, http.StatusInternalServerError, "Error scanning tasks", err)
			return
		}

		// Assume your MySQL TIMESTAMP is stored in UTC
        loc, _ := time.LoadLocation("Asia/Jakarta") // Your MySQL server timezone
        dateString := string(dateBytes)
        task.Date, err = time.ParseInLocation("2006-01-02 15:04:05", dateString, loc)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error parsing date", err)
            return
        }
        
        // Format the date string (hour-minutes-second - days-month-years)
        task.FormattedDate = task.Date.Format("15:04:05 || 02-01-2006")

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error.HandleError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var task Task
	err = database.DB.QueryRow("SELECT id, title, description, completed FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			error.HandleError(w, http.StatusNotFound, "Task not found", err)
		} else {
			error.HandleError(w, http.StatusInternalServerError, "Error fetching task", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// 1. Extract Task ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error.HandleError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	// 2. Decode JSON Request Body into a Task object
	var task *Task // Task can be nil
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		error.HandleError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// 3. Verify Task Existence
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		error.HandleError(w, http.StatusNotFound, "Task not found", err)
		return
	}

	// 4. Validate Task ID Match
	if task != nil && task.ID != int(id) {
		error.HandleError(w, http.StatusBadRequest, "Mismatched task ID", err)
		return
	}

	// 5. Determine Fields to Update (Partial Update Logic)
	updateFields := make(map[string]interface{})
	if task != nil {
		if task.Title != "" {
			updateFields["title"] = task.Title
		}
		if task.Description != "" {
			updateFields["description"] = task.Description
		}
	}
	updateFields["completed"] = task.Completed // Always update completed

	// 6. Construct Dynamic SQL Update Query
	var queryParts []string
	var queryArgs []interface{}
	for field, value := range updateFields {
		queryParts = append(queryParts, fmt.Sprintf("%s = ?", field))
		queryArgs = append(queryArgs, value)
	}
	queryArgs = append(queryArgs, id)

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(queryParts, ", "))

	// 7. Update Task in Database
	_, err = database.DB.Exec(query, queryArgs...)
	if err != nil {
		error.HandleError(w, http.StatusInternalServerError, "Failed to update task", err)
		return
	}

    var dateBytes []byte // To store raw date bytes
    row := database.DB.QueryRow("SELECT id, title, description, completed, date FROM tasks WHERE id = ?", id)
    err = row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &dateBytes)
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Failed to fetch updated task", err)
        return
    }

    // Parse the dateBytes into a time.Time value
    loc, _ := time.LoadLocation("Asia/Jakarta") // Your MySQL server timezone
    dateString := string(dateBytes)
    task.Date, err = time.ParseInLocation("2006-01-02 15:04:05", dateString, loc)
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error parsing date", err)
        return
    }
    
    // Format the date string (hour-minutes-second - days-month-years)
    task.FormattedDate = task.Date.Format("15:04:05 || 02-01-2006")

	// 9. Return Updated Task as JSON Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error.HandleError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	_, err = database.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		error.HandleError(w, http.StatusInternalServerError, "Failed to delete task", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
