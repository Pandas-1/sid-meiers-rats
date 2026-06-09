package controllers

import (
    "encoding/json"
    "net/http"
    "rats/models"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "time"
)


var jwtSecret = []byte("inthehistoryofjoeoverithasneverbeenthisjoeveridontevenknowifthisissecurebutafterenoughlettersithastobesecureatthispointlikewhat")

func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"` 
		Password string `json:"password"`
	}
	json.NewDecoder(r.body).Decode(&input)
	hash, err := bcrypt.GenerateFromPassword([](input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	err = models.CreateUser(input.Username, string(hash))
	if err != nil {
		http.Error(w, err , 400)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "registered successfully"})
}

