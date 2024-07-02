package work

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"backend/models"
	"backend/database"
)

type Work = models.Work

func CreateWork(w http.ResponseWriter, r *http.Request) {
	var work Work
	err := json.NewDecoder(r.Body).Decode(&work)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO works (title, description, completed) VALUES (?, ?, ?)", work.Title, work.Description, work.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastInsertID, _ := result.LastInsertId()
	work.ID = int(lastInsertID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(work)
}

func GetWorks(w http.ResponseWriter, r *http.Request) {
	var works []Work
	rows, err := database.DB.Query("SELECT id, title, description, completed FROM works")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var work Work
		err := rows.Scan(&work.ID, &work.Title, &work.Description, &work.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		works = append(works, work)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(works)
}

func GetWorkByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid work ID", http.StatusBadRequest)
		return
	}

	var work Work
	err = database.DB.QueryRow("SELECT id, title, description, completed FROM works WHERE id = ?", id).Scan(&work.ID, &work.Title, &work.Description, &work.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Work not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(work)
}

func UpdateWork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid work ID", http.StatusBadRequest)
		return
	}

	var work Work
	err = json.NewDecoder(r.Body).Decode(&work)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE works SET title = ?, description = ?, completed = ? WHERE id = ?", work.Title, work.Description, work.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteWork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid work ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM works WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
