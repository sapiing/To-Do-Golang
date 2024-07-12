package work

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"

	"backend/database"
	"backend/error"
	"backend/models"
)

type WorkLog = models.Work
type Task = models.Task

func GetWorkLogs(w http.ResponseWriter, r *http.Request) {
    var workLogs []models.Work
    rows, err := database.DB.Query(`
        SELECT wl.id, t.title, t.description, wl.date, wl.completed, wl.created_at, wl.updated_at
        FROM work_log_history wl
        JOIN tasks t ON wl.task_id = t.id 
        ORDER BY wl.date DESC
    `) // Updated query to get task title
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error querying work log history", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var workLog models.Work 
        var dateBytes, createdAtBytes, updatedAtBytes []byte 

        err := rows.Scan(&workLog.ID, &workLog.Title, &workLog.Description, &dateBytes, &workLog.Completed, &createdAtBytes, &updatedAtBytes)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error scanning work logs", err)
            return
        }

        loc, _ := time.LoadLocation("Asia/Jakarta") // Your MySQL server timezone
        
        // Parsing Time for Date
        dateString := string(dateBytes)

        if len(dateString) > 10 { // If it's a full timestamp, use this format
            workLog.Date, err = time.ParseInLocation("2006-01-02 15:04:05", dateString, loc)
        } else {  // If it's only a date, use this format
            workLog.Date, err = time.ParseInLocation("2006-01-02", dateString, loc)
        }
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error parsing date", err)
            return
        }
        
        // Parsing Time for CreatedAt
        createdAtString := string(createdAtBytes)
        workLog.CreatedAt, err = time.ParseInLocation("2006-01-02 15:04:05", createdAtString, loc)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error parsing created_at", err)
            return
        }
        
        // Parsing Time for UpdatedAt
        updatedAtString := string(updatedAtBytes)
        workLog.UpdatedAt, err = time.ParseInLocation("2006-01-02 15:04:05", updatedAtString, loc)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error parsing updated_at", err)
            return
        }
        // Format the date string (hour-minutes-second - days-month-years)
        workLog.FormattedDate = workLog.Date.Format("15:04:05 || 02-01-2006")

        workLogs = append(workLogs, workLog)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(workLogs)
}


func GetTodayWorkLogs(w http.ResponseWriter, r *http.Request) {
    today := time.Now().Format("2006-01-02")
    rows, err := database.DB.Query(`
        SELECT t.id, t.title, t.description, t.completed, wl.completed AS daily_completed
        FROM tasks t
        LEFT JOIN work_log wl ON t.id = wl.task_id AND wl.date = ?
    `, today)
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error querying work logs", err)
        return
    }
    defer rows.Close()

    var workLogs []WorkLog
    for rows.Next() {
        var workLog WorkLog
        var dailyCompleted sql.NullBool
        err := rows.Scan(&workLog.ID, &workLog.Title, &workLog.Description, &workLog.Completed, &dailyCompleted) // Include Title and Description
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error scanning work logs", err)
            return
        }

        workLog.DailyCompleted = dailyCompleted.Bool

        workLogs = append(workLogs, workLog)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(workLogs)
}

// AssignDailyTasks assigns daily tasks for the current day, moving old tasks to history.
func AssignDailyTasks(w http.ResponseWriter, r *http.Request) {
    log.Println("AssignDailyTasks function called")
    today := time.Now().Format("2006-01-02")

    // 1. Pindahkan semua work log ke history (tanpa filter tanggal)
    _, err := database.DB.Exec("INSERT INTO work_log_history (task_id, date, completed, created_at) SELECT task_id, date, completed, created_at FROM work_log")
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error moving work logs to history", err)
        return
    }

    // 2. Hapus semua work log (tanpa filter tanggal)
    _, err = database.DB.Exec("DELETE FROM work_log")
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error deleting work logs", err)
        return
    }

    // 3. Ambil daftar tugas yang belum selesai
    var tasks []Task
    rows, err := database.DB.Query("SELECT * FROM tasks WHERE completed = 0")
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Error querying tasks", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var task Task
        var dateBytes []byte
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &dateBytes)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error scanning tasks", err) 
            return
        }
        tasks = append(tasks, task)
    }

    // 4. Masukkan tugas harian baru (dengan tanggal hari ini)
    for _, task := range tasks {
        _, err := database.DB.Exec("INSERT IGNORE INTO work_log (task_id, date, completed) VALUES (?, ?, ?)", task.ID, today, 0)
        if err != nil {
            error.HandleError(w, http.StatusInternalServerError, "Error inserting work log", err)
            return
        }
    }
    
    // Success response - if you need it
    w.WriteHeader(http.StatusCreated) 
    json.NewEncoder(w).Encode(map[string]string{"message": "Daily tasks assigned"})
}




func UpdateWorkLog(w http.ResponseWriter, r *http.Request) {
    var requestData struct { 
        TaskID    int    `json:"taskId"`
        Completed bool   `json:"completed"`
        Date      string `json:"date"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        error.HandleError(w, http.StatusBadRequest, "Invalid request body", err)
        return
    }
    
    // Update data di database (dengan kondisi tanggal)
    _, err := database.DB.Exec("UPDATE work_log SET completed = ? WHERE task_id = ? AND date = ?", requestData.Completed, requestData.TaskID, requestData.Date)
    if err != nil {
        error.HandleError(w, http.StatusInternalServerError, "Failed to update work log", err)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Work log updated"})
}



