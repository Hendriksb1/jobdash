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
	Model            // embedded so all methods of model are availiable on server
	JobIdRelation    map[string]int
	ResultIdRelation map[string]int
}

type Model struct {
	*sql.DB
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
		fmt.Printf("Error opening database: %v", err)
	}

	err = s.DB.Ping()
	if err != nil {
		fmt.Printf("Error connecting to database: %v", err)
	}

	// laod maps
	s.JobIdRelation, err = s.LoadIdNameRelation("job_types")
	if err != nil {
		fmt.Printf("failed to load job types: %v", err.Error())
	}

	s.ResultIdRelation, err = s.LoadIdNameRelation("results")
	if err != nil {
		fmt.Printf("failed to load result types: %v", err.Error())
	}

	// update ghost status
	err := s.UpdateGhostedStatus()
	if err != nil {
		fmt.Printf("failed to update ghost status: %v", err.Error())
	}

	// update weekly table
	t, err := s.GetLastResetDate()
	if err != nil {
		fmt.Printf("failed to get last reset date", err.Error())
	}

	// continue here with figuring out the logic

	// today is monday
	// next is today
	fmt.Println(t)


	mux := mux.NewRouter()

	// openings
	mux.HandleFunc("/addOpening", s.AddOpening).Methods(http.MethodPost)
	mux.HandleFunc("/getAllOpenings/{id}", s.GetAllOpenings).Methods(http.MethodGet)
	mux.HandleFunc("/getOpening", s.GetOpening).Methods(http.MethodGet)
	mux.HandleFunc("/getOpeningsByResult", s.GetOpeningsByResult)
	mux.HandleFunc("/deleteOpening/{id}", s.DeleteOpening).Methods(http.MethodDelete)
	mux.HandleFunc("/updateOpening/{id}", s.UpdateOpening).Methods(http.MethodPut)

	// types
	mux.HandleFunc("/getAllResultTypes", s.GetAllResultTypes).Methods(http.MethodGet)
	mux.HandleFunc("/getAllJobTypes", s.GetAllJobTypes).Methods(http.MethodGet)

	// aggregations
	mux.HandleFunc("/getOpeningsByJob", s.GetOpeningsByJob)
	mux.HandleFunc("/getCountThisWeek", s.GetCountThisWeek)
	mux.HandleFunc("/getCountsPerWeek", s.GetCountsPerWeek)

	// job type
	mux.HandleFunc("/addJobType", s.AddJobType)

	// user
	mux.HandleFunc("/registerUser", s.RegisterUser)
	mux.HandleFunc("/unRegisterUser", s.UnRegisterUser)
	mux.HandleFunc("/changeUser", s.ChangeUser)
	mux.HandleFunc("/getUser", s.GetUser)
	mux.HandleFunc("/getUserByEmail", s.GetUserByEmail)

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

	r, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("updated %v rows due to ghosting", r)

	return nil
}
