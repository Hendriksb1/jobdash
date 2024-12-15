package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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
