package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var users []User
var sessions = make(map[string]bool)

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/notes", handleNotes)
	http.ListenAndServe(":8080", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	sid := user.Email
	sessions[sid] = true

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"sid": sid,
	})
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		listNotes(w, r)
	} else if r.Method == http.MethodPost {
		createNote(w, r)
	} else if r.Method == http.MethodDelete {
		deleteNote(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sid := r.FormValue("sid")
	if !isLoggedIn(sid) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notes := []Note{
		{ID: 1, Note: "New Note first"},
		{ID: 2, Note: "New Note second"},
		{ID: 3, Note: "New Note third"},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Note{
		"notes": notes,
	})
}

func createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sid := r.FormValue("sid")
	if !isLoggedIn(sid) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	note.ID = 123

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]uint32{
		"id": note.ID,
	})
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sid := r.FormValue("sid")
	if !isLoggedIn(sid) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isLoggedIn(sid string) bool {
	return sessions[sid]
}
