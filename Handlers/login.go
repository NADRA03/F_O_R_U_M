package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    "github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")

        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            password := r.FormValue("password")

            var id int
            var storedPassword string
            err := db.QueryRow("SELECT id, password FROM user WHERE username = ?", username).Scan(&id, &storedPassword)
            if err != nil || password != storedPassword {
                http.Error(w, "Invalid username or password", http.StatusUnauthorized)
                return
            }

            session.Values["id"] = id
            session.Values["authenticated"] = true
            session.Values["username"] = username
            session.Save(r, w)

            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        tmpl, err := template.ParseFiles("HTML/Login.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
    }
}


