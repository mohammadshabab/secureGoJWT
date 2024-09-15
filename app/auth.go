package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"securegojwt/models"
	"securegojwt/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Define a custom key type to avoid collisions
type ContextKey string

const usr ContextKey = "user"

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//These endpoints do not requried authentication
		noAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path

		//Check if req does not need authentication
		for _, value := range noAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		//Get token from header
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = utils.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		fmt.Println("err=========== ", err)
		if err != nil {
			response = utils.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}
		if !token.Valid {
			response = utils.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		//set the caller to the user retrieved from the parsed token
		fmt.Printf("User %d", tk.UserId)

		ctx := context.WithValue(r.Context(), usr, tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
