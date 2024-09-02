package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func main() {
	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	http.HandleFunc("/users", GetUsersHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		http.Error(w, `{"message": "Failed to query users", "status_code": 500}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			http.Error(w, `{"message": "Failed to scan user", "status_code": 500}`, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := map[string]interface{}{
		"message":     "success",
		"status_code": 200,
		"data":        users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
