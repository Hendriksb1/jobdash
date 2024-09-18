package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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
		fmt.Println(err)
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var counts []WeeklyCount
	for rows.Next() {
		var count WeeklyCount
		if err := rows.Scan(&count.WeekEndDate, &count.EntryCount); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		counts = append(counts, count)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(counts); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
