package handlers

import (
	"encoding/json"
	"net/http"
	"video-streaming/db"
	"video-streaming/models"
)

const (
	DATABASE = "video-stream"
	URI      = "mongodb+srv://j:rootroot@cluster0.rj0tg.mongodb.net/"
)

func InsertUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application-json")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.InitDatabase(DATABASE, URI).InsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response struct {
		Message string `json:"message"`
		Error   bool   `json:"error"`
		ID      string `json:"_id"`
	}
	response.Message = "User successfuly inserted"
	response.Error = false
	response.ID = id.Hex()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.InitDatabase(DATABASE, URI).GetUser(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Password != u.Password {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	var response struct {
		Message string       `json:"message"`
		Error   bool         `json:"error"`
		User    *models.User `json:"user"`
	}
	response.Message = "User successfuly inserted"
	response.Error = false
	user.Password = ""
	response.User = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
