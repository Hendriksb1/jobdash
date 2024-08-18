package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// | CREATE TABLE user(id INTEGER primary key, name text, email text)
type User struct {
	Id    int32  `json:"id"` // primary key
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) InsertUser(u *User) error {

	q := `
		INSERT INTO users (name, email) values(?, ?);
	`
	res, err := s.DB.Exec(q, u.Name, u.Email)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

// TODO continue with delete user and update users

func (s *Server) DeleteUser(id int32) error {

	q := `
		DELETE FROM users WHERE id = ?;
	`

	res, err := s.DB.Exec(q, id)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func (s *Server) UpdateUser(u *User) error {
	q := `
		UPDATE users SET (name, email) values(?, ?) WHERE id = ?;
	`

	res, err := s.DB.Exec(q, u.Name, u.Email, u.Id)
	if err != nil {
		return err
	}

	fmt.Println(res)
}
