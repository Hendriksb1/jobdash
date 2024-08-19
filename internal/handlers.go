package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Opening represents a job opening
type Opening struct {
	ID              int    `json:"id"`
	Firm            string `json:"firm"`
	TypeJob         string `json:"type_job"`
	Result          string `json:"result"`
	ApplicationDate string `json:"application_date,omitempty"`
	URL             string `json:"url"`
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

	if  _, ok := s.JobIdRelation[opening.TypeJob]; !ok {
		// update job relation map
		s.JobIdRelation, err = s.LoadIdNameRelation("job_types")
		if err != nil {
			print("failed to load job types: %v", err.Error())
		}
	}

	query := "INSERT INTO openings (firm, type_job, result, url) VALUES (?, ?, ?, ?)"
	result, err := s.DB.Exec(query, opening.Firm, s.JobIdRelation[opening.TypeJob], s.ResultIdRelation[opening.Result], opening.URL)
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

func (s *Server) GetAllOpenings(w http.ResponseWriter, r *http.Request) {

	fmt.Println("attempting to get all openings")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := s.DB.Query("SELECT openings.id, firm, type_name, result_name, application_date, url FROM openings JOIN results ON openings.result = results.id JOIN job_types ON openings.type_job = job_types.id")
	if err != nil {
		http.Error(w, "Failed to retrieve openings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var openings []Opening
	for rows.Next() {
		var opening Opening
		if err := rows.Scan(&opening.ID, &opening.Firm, &opening.TypeJob, &opening.Result, &opening.ApplicationDate, &opening.URL); err != nil {
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

func (s *Server) GetAllJobTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := s.DB.Query("SELECT type_name FROM job_types")
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
		SELECT results.result_name, COUNT(openings.id) as count_total 
		FROM openings 
		JOIN results ON results.id = openings.result 
		GROUP BY results.result_name
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

	query := `
		SELECT COUNT(*)
        FROM openings
        WHERE strftime('%Y-%W', application_date) = strftime('%Y-%W', 'now');
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch openings count for this week", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
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

type WeeklyCount struct {
	WeekEndDate string `json:"week_end_date"`
	EntryCount  int    `json:"entry_count"`
}

func (s *Server) GetCountsPerWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT week_end_date, entry_count FROM weekly_counts;
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var counts []WeeklyCount
	for rows.Next() {
		var count WeeklyCount
		if err := rows.Scan(&count.WeekEndDate, &count.EntryCount); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		counts = append(counts, count)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(counts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
