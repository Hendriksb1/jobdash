package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Opening represents a job opening
type Opening struct {
	ID              int    `json:"id"`
	Firm            string `json:"firm"`
	TypeJob         string `json:"type_job"`
	Result          string `json:"result"`
	ApplicationDate string `json:"application_date,omitempty"`
	URL             string `json:"url"`
	UserID          int    `json:"user_id"`
}

func (s *Server) AddOpening(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var opening Opening
	err := json.NewDecoder(r.Body).Decode(&opening)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, ok := s.JobIdRelation[opening.TypeJob]; !ok {
		// update job relation map
		s.JobIdRelation, err = s.LoadIdNameRelation("job_types")
		if err != nil {
			print("failed to load job types: %v", err.Error())
		}
	}

	query := "INSERT INTO openings (firm, type_job, result, url, user_id) VALUES (?, ?, ?, ?, ?)"
	result, err := s.DB.Exec(query, opening.Firm, s.JobIdRelation[opening.TypeJob], s.ResultIdRelation[opening.Result], opening.URL, opening.UserID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to insert opening", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	opening.ID = int(id)
	// opening.ApplicationDate = "date('now')"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opening)

	fmt.Println("added ", opening)
}

func (s *Server) GetOpening(w http.ResponseWriter, r *http.Request) {
	fmt.Println("attempting to get an opening")

	if r.Method != http.MethodGet {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	userIdStr := r.URL.Query().Get("userId")
	jobIdStr := r.URL.Query().Get("jobId")

	fmt.Println(userIdStr)
	fmt.Println(jobIdStr)

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	q := `
	SELECT openings.id, firm, type_name, result_name, application_date, url
	FROM openings
	JOIN results ON openings.result = results.id
	JOIN job_types ON openings.type_job = job_types.id
	WHERE openings.user_id = ? AND openings.id = ?;
	`

	var opening struct {
		ID              int       `json:"id"`
		Firm            string    `json:"firm"`
		TypeName        string    `json:"type_name"`
		ResultName      string    `json:"result_name"`
		ApplicationDate time.Time `json:"application_date"`
		URL             string    `json:"url"`
	}

	err = s.DB.QueryRow(q, userId, jobId).Scan(&opening.ID, &opening.Firm, &opening.TypeName, &opening.ResultName, &opening.ApplicationDate, &opening.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Opening not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the opening struct to JSON and write to response
	if err := json.NewEncoder(w).Encode(opening); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetAllOpenings is used for the list of openings
func (s *Server) GetAllOpenings(w http.ResponseWriter, r *http.Request) {

	fmt.Println("attempting to get all openings")

	if r.Method != http.MethodGet {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	q := `
	SELECT openings.id, firm, type_name, result_name, application_date, url
	FROM openings
	JOIN results ON openings.result = results.id
	JOIN job_types ON openings.type_job = job_types.id
	WHERE openings.user_id = ?;
	`

	fmt.Println(q)
	fmt.Println(userId)

	rows, err := s.DB.Query(q, userId)
	if err != nil {
		fmt.Println("Failed to retrieve openings")
		http.Error(w, "Failed to retrieve openings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var openings []Opening
	for rows.Next() {
		var opening Opening
		if err := rows.Scan(&opening.ID, &opening.Firm, &opening.TypeJob, &opening.Result, &opening.ApplicationDate, &opening.URL); err != nil {
			fmt.Println("Failed to scan opening")
			http.Error(w, "Failed to scan opening", http.StatusInternalServerError)
			return
		}
		openings = append(openings, opening)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error occurred during row iteration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openings)
}

func (s *Server) DeleteOpening(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid opening ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM openings WHERE id = ?"
	result, err := s.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete opening", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows deleted", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Opening deleted successfully"))
}

func (s *Server) UpdateOpening(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var opening Opening
	err := json.NewDecoder(r.Body).Decode(&opening)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "UPDATE openings SET firm = ?, type_job = ?, result = ?, application_date = ?, url = ? WHERE id = ?"
	_, err = s.DB.Exec(query, opening.Firm, s.JobIdRelation[opening.TypeJob], s.ResultIdRelation[opening.Result], opening.ApplicationDate, opening.URL, id)
	if err != nil {
		http.Error(w, "Failed to update opening", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(opening)
}
