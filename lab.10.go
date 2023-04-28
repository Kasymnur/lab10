package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int
    Name  string
    Email string
}

var db *sql.DB
var tmpl = template.Must(template.ParseGlob("templates/*"))

func main() {
   
    var err error
    db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/database_name")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

   
    http.HandleFunc("/", index)
    http.HandleFunc("/create", create)
    http.HandleFunc("/update", update)
    http.HandleFunc("/delete", delete)
    log.Fatal(http.ListenAndServe(":8080", nil))
}


func create(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        name := r.FormValue("name")
        email := r.FormValue("email")

        
        db.Exec("INSERT INTO users(name, email) VALUES(?, ?)", name, email)

      
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    
    tmpl.ExecuteTemplate(w, "create.html", nil)
}


func index(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            log.Fatal(err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    tmpl.ExecuteTemplate(w, "index.html", users)
}


func update(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id := r.FormValue("id")
        name := r.FormValue("name")
        email := r.FormValue("email")

        db.Exec("UPDATE users SET name=?, email=? WHERE id=?", name, email, id)

     
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

   
    id := r.FormValue("id")
    row := db.QueryRow("SELECT * FROM users WHERE id=?", id)
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        log.Fatal(err)
    }
    tmpl.ExecuteTemplate(w, "update.html", user)
}


func delete(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")

    db.Exec("DELETE FROM users WHERE id=?", id)

    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}
