package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Result represents a result type
type Result struct {
	ResultName string `json:"result_name"`
}

type Result2 struct {
	TypeName string `json:"type_name"`
}

func (s *Server) GetAllResultTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := s.DB.Query("SELECT result_name FROM results")
	if err != nil {
		http.Error(w, "Failed to retrieve result types", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var result Result
		if err := rows.Scan(&result.ResultName); err != nil {
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

// ResultOverview represents the structure for result overview
type ResultOverview struct {
	ResultName string `json:"result_name"`
	CountTotal int    `json:"count_total"`
}

func (s *Server) GetOpeningsByResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT results.result_name, COUNT(openings.id) as total_count 
		FROM openings 
		JOIN results ON results.id = openings.result 
		GROUP BY results.result_name
		ORDER BY total_count DESC;
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch openings by result", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var resultOverview []ResultOverview

	for rows.Next() {
		var overview ResultOverview
		err := rows.Scan(&overview.ResultName, &overview.CountTotal)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		resultOverview = append(resultOverview, overview)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultOverview)
}

type JobOverview struct {
	Type_name  string `json:"type_name"`
	CountTotal int    `json:"count_total"`
}

// SELECT type_name, count(*) as total_count FROM openings JOIN job_types as j ON j.id = openings.type_job GROUP BY type_job
func (s *Server) GetOpeningsByJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT type_name,
		count(*) as total_count
		FROM openings
		JOIN job_types as j ON j.id = openings.type_job 
		GROUP BY type_job
		ORDER BY total_count DESC;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch openings by job type", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var resultOverview []JobOverview

	for rows.Next() {
		var overview JobOverview
		err := rows.Scan(&overview.Type_name, &overview.CountTotal)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		resultOverview = append(resultOverview, overview)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultOverview)
}

// SELECT COUNT(*)FROM openings WHERE strftime('%Y-%W', application_date) = strftime('%Y-%W', 'now');
func (s *Server) GetCountThisWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// result 6 is "want to apply" and this will not be counted
	query := `
		SELECT COUNT(*)
    	FROM openings
    	WHERE strftime('%Y-%W', application_date) = strftime('%Y-%W', 'now') AND result != 6;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch openings count for this week", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	resultOverview := struct {
		Count int `json:"count"`
	}{
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultOverview)
}
