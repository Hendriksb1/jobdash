package internal

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (s *Server) GetLastResetDate() (time.Time, error) {
	var eventDate time.Time

	q1 := `SELECT event_date
		FROM last_reset 
		WHERE id = 1;`

	// Use QueryRow to retrieve a single row
	err := s.DB.QueryRow(q1).Scan(&eventDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, fmt.Errorf("no reset date found")
		}
		return time.Time{}, fmt.Errorf("failed to retrieve last reset date: %w", err)
	}

	return eventDate, nil
}

// ResetWeeklyCounts is setting the weekly count back to 0
// also it stores the last reset date...so its only resetted once a week
func (s *Server) ResetWeeklyCounts() error {

	// Begin the transaction
	tx, err := s.DB.Begin()
	if err != nil {
		fmt.Printf("could not begin transaction: %v", err)
		return err
	}

	// Defer the rollback to ensure it's called in case of an error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	q1 := `UPDATE users 
		  SET weekly_applications_count = 0;`

	_, err = tx.Exec(q1)
	if err != nil {
		return err
	}

	// TODO: continue here, need to create table and store the date
	q2 := `Insert into last_reset (event_date) values (?);`

	_, err = tx.Exec(q2, time.Now())
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("successfully reset weekly counts")
	return nil
}
