package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type LoggedInUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CookieData struct {
	LoggedInUser LoggedInUser `json:"loggedInUser"`
}

func GetUserIdFromCookie(r *http.Request, w http.ResponseWriter) int {
	// Retrieve the cookie
	cookie, err := r.Cookie("jobDash")
	if err != nil {
		http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
		return 0
	}

	// URL decode the cookie value
	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to decode cookie", http.StatusBadRequest)
		return 0
	}

	// Unmarshal the JSON from the decoded value
	var cookieData CookieData
	err = json.Unmarshal([]byte(decodedValue), &cookieData)
	if err != nil {
		http.Error(w, "Failed to parse cookie JSON", http.StatusBadRequest)
		return 0
	}

	// Return the user ID
	userID := cookieData.LoggedInUser.ID
	return userID
}
