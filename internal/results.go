package internal

import (
	"encoding/json"
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

// ResultOverview represents the structure for result overview
type ResultOverview struct {
	ResultName string `json:"result_name"`
	CountTotal int    `json:"count_total"`
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

func (s *Server) GetOpeningsByResult(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a GET request
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the user cookie (adjust the name to match your cookie name)
	// Replace "user_session" with the actual cookie name
	// URL-decode the cookie value
	// Deserialize the JSON data
	// Extract the logged-in user's name from the deserialized data
	userID := GetUserIdFromCookie(r, w)
	if userID == 0 {
		return
	}

	// Modify the query to filter by user-specific results
	query := `
		SELECT results.result_name, COUNT(openings.id) as total_count 
		FROM openings 
		JOIN results ON results.id = openings.result 
		WHERE openings.user_id = ?  -- Filter by user_id
		GROUP BY results.result_name
		ORDER BY total_count DESC;
	`

	// Execute the query with the userID as a parameter
	rows, err := s.DB.Query(query, userID)
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

	// Check for any errors that may have occurred during iteration
	err = rows.Err()
	if err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the result overview data as JSON and write it to the response
	json.NewEncoder(w).Encode(resultOverview)
}
