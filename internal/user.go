package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	// "net/url"
)

// | CREATE TABLE user(id INTEGER primary key, name text, email text)
type User struct {
	Id    int32  `json:"id"` // primary key
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("attempting to get user")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if !r.URL.Query().Has("id") {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	resUser, err := s.SelectUser(int32(id))
	if err != nil {
		fmt.Println("failed to select user: ", err.Error())
		http.Error(w, "failed to select user", http.StatusInternalServerError)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resUser)
}

func (m *Model) SelectUser(id int32) (*User, error) {
	q := `
		SELECT id, name, email FROM users WHERE id = ?;
	`
	rows, err := m.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Check if there is a result
	if !rows.Next() {
		return nil, errors.New("no user found")
	}

	// Create a User instance
	user := &User{}

	// Assuming your User struct has fields: ID, Name, and Email
	if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
		return nil, err
	}

	return user, nil
}


func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.InsertUser(&u)
	if err != nil {
		fmt.Println("failed to store user: ", err.Error())
		http.Error(w, "failed to store user", http.StatusInternalServerError)
		return
	}
}

func (m *Model) InsertUser(u *User) error {

	q := `
		INSERT INTO users (name, email) values(?, ?);
	`
	res, err := m.Exec(q, u.Name, u.Email)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

type UnRegisterUserRequest struct {
	Id string `json:"id"`
}

func (s *Server) UnRegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req UnRegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.Id)
	if err != nil {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	err = s.DeleteUser(int32(id))
	if err != nil {
		fmt.Println("failed to store user: ", err.Error())
		http.Error(w, "failed to un-register user", http.StatusInternalServerError)
		return
	}
}

// TODO continue with delete user and update users

func (m *Model) DeleteUser(id int32) error {

	q := `
		DELETE FROM users WHERE id = ?;
	`

	res, err := m.Exec(q, id)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func (s *Server) ChangeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.UpdateUser(&u)
	if err != nil {
		fmt.Println("failed to change user: ", err.Error())
		http.Error(w, "failed to change user", http.StatusInternalServerError)
		return
	}
}

func (m *Model) UpdateUser(u *User) error {
	q := `
		UPDATE users SET (name, email) values(?, ?) WHERE id = ?;
	`

	res, err := m.Exec(q, u.Name, u.Email, u.Id)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
