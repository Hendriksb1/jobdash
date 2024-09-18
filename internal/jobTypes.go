package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type JobType struct {
	Name string `json:"name"`
}

// AddJobType is the handler function for adding jobs
func (s *Server) AddJobType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var j JobType
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.InsertJobType(j.Name)
	if err != nil {
		fmt.Println("failed to store job name. Error code: ", err.Error())
		return
	}
}

// InsertJobType is the database method for adding a new job type
func (s *Server) InsertJobType(jobName string) error {

	q := `
		INSERT INTO job_types (type_name) values(?);
	`
	_, err := s.DB.Exec(q, jobName)
	if err != nil {
		return fmt.Errorf("%i", http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) GetAllJobTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT type_name 
		FROM job_types
		ORDER BY LOWER(type_name) DESC;	
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to retrieve job types", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Result2
	for rows.Next() {
		var result Result2
		if err := rows.Scan(&result.TypeName); err != nil {
			http.Error(w, "Failed to scan result type", http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error occurred during row iteration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
