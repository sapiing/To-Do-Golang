package models

import "time"

type Work struct {
    ID        int       `json:"id"`
    TaskID    int       `json:"task_id"`
    Date      time.Time `json:"date"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    FormattedDate string `json:"formatted_date"`
    Title string `json:"title"`          // Add this field
    Description string `json:"description"`  // Add this field
    DailyCompleted bool      `json:"daily_completed"`
}

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    Date        time.Time `json:"date"`
    FormattedDate string `json:"formattedDate"`
}