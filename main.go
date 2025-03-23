package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var tmpl *template.Template

type User struct {
	ID   int
	Name string
}

func main() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	time.Sleep(5 * time.Second)

	// Connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table & insert fake data
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT);
					  INSERT INTO users (name) VALUES
					   ('Alice'), ('Bob'), ('Charlie'),
					   ('David'), ('Eve'), ('Frank') ON CONFLICT DO NOTHING;
					  `)
	if err != nil {
		log.Fatal(err)
	}

	// Parse template
	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	// Serve HTTP
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name); err != nil {
				log.Println(err)
			}
			users = append(users, u)
		}
		tmpl.Execute(w, users)
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
