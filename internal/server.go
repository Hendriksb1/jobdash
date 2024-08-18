package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var err error

type Server struct {
	DB               *sql.DB
	JobIdRelation    map[string]int
	ResultIdRelation map[string]int
}

func (s *Server) Init() {

	port := 8080
	if isPortInUse(port) {
		log.Fatalf("Port %d is already in use", port)
	}

	fmt.Println("starting server")

	defer s.DB.Close()

	// Initialize the CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	s.DB, err = sql.Open("sqlite3", "./findajob.db")
	if err != nil {
		print("Error opening database: %v", err)
	}

	err = s.DB.Ping()
	if err != nil {
		print("Error connecting to database: %v", err)
	}

	// laod maps
	s.JobIdRelation, err = s.LoadIdNameRelation("job_types")
	if err != nil {
		print("failed to load job types: %v", err.Error())
	}

	s.ResultIdRelation, err = s.LoadIdNameRelation("results")
	if err != nil {
		print("failed to load result types: %v", err.Error())
	}

	// update ghost status
	err := s.UpdateGhostedStatus()
	if err != nil {
		print("failed to update ghost status: %v", err.Error())
	}

	// update weekly table
	err = s.UpdateWeeklyCounts()
	if err != nil {
		print("failed to update weekly counts: %v", err.Error())
	}

	// set up handlers
	mux := mux.NewRouter()
	mux.HandleFunc("/addOpening", s.AddOpening).Methods(http.MethodPost)
	mux.HandleFunc("/getAllOpenings", s.GetAllOpenings).Methods(http.MethodGet)
	mux.HandleFunc("/getAllResultTypes", s.GetAllResultTypes).Methods(http.MethodGet)
	mux.HandleFunc("/getAllJobTypes", s.GetAllJobTypes).Methods(http.MethodGet)
	mux.HandleFunc("/deleteOpening/{id}", s.DeleteOpening).Methods(http.MethodDelete)
	mux.HandleFunc("/updateOpening/{id}", s.UpdateOpening).Methods(http.MethodPut)
	mux.HandleFunc("/getOpeningsByResult", s.GetOpeningsByResult)
	mux.HandleFunc("/getOpeningsByJob", s.GetOpeningsByJob)
	mux.HandleFunc("/getCountThisWeek", s.GetCountThisWeek)
	mux.HandleFunc("/getCountsPerWeek", s.GetCountsPerWeek)
	mux.HandleFunc("/addJobType", s.AddJobType)

	handler := c.Handler(mux)
	fmt.Printf("Server is running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))

	// http.HandleFunc("/addOpening", s.AddOpening)
	// fmt.Println("Server is running on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func isPortInUse(port int) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// UpdateGhostedStatus will change the ghost result status
// later this should be a go routine
func (s *Server) UpdateGhostedStatus() error {
	q := `UPDATE openings 
		  SET result = 5
		  WHERE result = 1 AND application_date < date('now', '-21 days');`

	res, err := s.DB.Exec(q)
	if err != nil {
		return err
	}

	r, err :=  res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("updated %v rows due to ghosting", r)

	return nil
}

func (s *Server) UpdateWeeklyCounts () error {
	q := `DELETE FROM weekly_counts;
        	 INSERT INTO weekly_counts (week_end_date, entry_count)
             SELECT
                 DATE(application_date, 'weekday 0', '+6 days') AS week_end_date,
                 COUNT(*) AS entry_count
             FROM
                 openings
             GROUP BY
                 week_end_date;`
	res, err := s.DB.Exec(q)
	if err != nil {
		return err
	}

	r, err :=  res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("updated %v rows due to weekly counting", r)

	return nil
}