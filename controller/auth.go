package controller

import (
    "encoding/json"
    "net/http"
    "rats/models"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "time"
	"context"
)

type contextKey string
const UserIDKey contextKey = "userID"
var jwtSecret = []byte("inthehistoryofjoeoverithasneverbeenthisjoeveridontevenknowifthisissecurebutafterenoughlettersithastobesecureatthispointlikewhat")

func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"` 
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	err = models.CreateUser(input.Username, string(hash))
	if err != nil {
		http.Error(w, err.Error() , 400)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	user, err := models.GetUserByUsername(input.Username)
	if err != nil {
		http.Error(w,err.Error(), 404)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		http.Error(w, err.Error(),401)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"exp": time.Now().Add(4*time.Hour).Unix(),
	})
	tokenString , err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "No token provided", 401)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", 401)
            return
        }
		claims := token.Claims.(jwt.MapClaims)
        userID := int(claims["userID"].(float64))
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}