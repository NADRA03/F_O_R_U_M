package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
)

func SignupHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            password := r.FormValue("password")
            email := r.FormValue("email")
            image := r.FormValue("image") // This is optional; adapt based on your needs

            // Check if username already exists
            var existingUser string
            err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&existingUser)
            if err != nil && err != sql.ErrNoRows {
                http.Error(w, "Database error", http.StatusInternalServerError)
                return
            }

            if existingUser != "" {
                http.Error(w, "Username already taken", http.StatusBadRequest)
                return
            }

            // Insert the new user into the database
            _, err = db.Exec("INSERT INTO user (username, password, email, image) VALUES (?, ?, ?, ?)", username, password, email, image)
            if err != nil {
                http.Error(w, "Database error", http.StatusInternalServerError)
                return
            }

            // Redirect to login page after successful signup
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        tmpl, err := template.ParseFiles("HTML/Signup.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
    }
}
