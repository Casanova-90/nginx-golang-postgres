package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
    // Configuración directa de la conexión a la base de datos
    host := "YOUR_CLOUD_SQL_PUBLIC_IP"
    port := "5432"
    user := "YOUR_DB_USER"
    password, err := ioutil.ReadFile("/run/secrets/db-password")
    if err != nil {
        return nil, err
    }
    dbname := "example"

    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, string(password), host, port, dbname)
    return sql.Open("postgres", connStr)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        w.WriteHeader(500)
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT title FROM blog")
    if err != nil {
        w.WriteHeader(500)
        return
    }
    var titles []string
    for rows.Next() {
        var title string
        err = rows.Scan(&title)
        titles = append(titles, title)
    }
    json.NewEncoder(w).Encode(titles)
}

func singleBlogHandler(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        w.WriteHeader(500)
        return
    }
    defer db.Close()

    params := mux.Vars(r)
    id := params["id"]

    var title string
    err = db.QueryRow("SELECT title FROM blog WHERE id = $1", id).Scan(&title)
    if err != nil {
        w.WriteHeader(500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"title": title})
}

func createBlogHandler(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        w.WriteHeader(500)
        return
    }
    defer db.Close()

    var newBlog struct {
        Title string `json:"title"`
    }
    err = json.NewDecoder(r.Body).Decode(&newBlog)
    if err != nil {
        w.WriteHeader(400)
        return
    }

    _, err = db.Exec("INSERT INTO blog (title) VALUES ($1)", newBlog.Title)
    if err != nil {
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(201)
}

func prepare() error {
    db, err := connect()
    if err != nil {
        return err
    }
    defer db.Close()

    for i := 0; i < 60; i++ {
        if err := db.Ping(); err == nil {
            break
        }
        time.Sleep(time.Second)
    }

    if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
        return err
    }

    if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id SERIAL, title VARCHAR)"); err != nil {
        return err
    }

    for i := 0; i < 5; i++ {
        if _, err := db.Exec("INSERT INTO blog (title) VALUES ($1);", fmt.Sprintf("Blog post #%d", i)); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    log.Print("Prepare db...")
    if err := prepare(); err != nil {
        log.Fatal(err)
    }

    log.Print("Listening on port 8000")
    r := mux.NewRouter()
    r.HandleFunc("/blogs", blogHandler).Methods("GET")
    r.HandleFunc("/blogs/{id}", singleBlogHandler).Methods("GET")
    r.HandleFunc("/blogs", createBlogHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}