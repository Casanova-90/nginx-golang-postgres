package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	// Get the environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// If any of the environment variables are missing, return an error
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Construct the connection string using the environment variables
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	return sql.Open("postgres", connStr) 
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error connecting to the database: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT title FROM blog")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error querying the database: %v", err)
		return
	}

	var titles []string
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error scanning row: %v", err)
			return
		}
		titles = append(titles, title)
	}

	json.NewEncoder(w).Encode(titles)
}

func main() {
	log.Print("Prepare db...")

	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", blogHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Retry until the database is available
	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// Drop the table if it exists and create a new one
	if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id SERIAL, title VARCHAR)"); err != nil {
		return err
	}

	// Insert some sample blog posts
	for i := 0; i < 5; i++ {
		if _, err := db.Exec("INSERT INTO blog (title) VALUES ($1);", fmt.Sprintf("Blog post #%d", i)); err != nil {
			return err
		}
	}
	return nil
}