package Forum

import (
	"database/sql"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	// "github.com/gorilla/sessions"
)

func EditProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")
        userID, ok := session.Values["id"].(int)
		// fmt.Println("userID: ", userID)
        if !ok || session.Values["authenticated"] != true {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            email := r.FormValue("email")
            password := r.FormValue("password")
            image := r.FormValue("image")

            var existingID int
            err := db.QueryRow("SELECT id FROM user WHERE id = ?", userID).Scan(&existingID)
            if err != nil {
                if err == sql.ErrNoRows {
                    http.Error(w, "User not found", http.StatusNotFound)
                } else {
                    http.Error(w, "Database error", http.StatusInternalServerError)
                }
                log.Println("User ID not found or other error:", userID, err)
                return
            }

            _, err = db.Exec("UPDATE user SET username = ?, email = ?, password = ?, image = ? WHERE id = ?",
                username, email, password, image, userID)
            if err != nil {
                http.Error(w, "Database error", http.StatusInternalServerError)
                log.Println("Error updating user:", err)
                return
            }

            session.Values["username"] = username
            session.Values["profileImage"] = image
            session.Save(r, w)

            http.Redirect(w, r, "/settings", http.StatusSeeOther)
            return
        }

        var username, email, image string
        err := db.QueryRow("SELECT username, email, image FROM user WHERE id = ?", userID).Scan(&username, &email, &image)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            log.Println("Error retrieving user information:", err)
            return
        }

        tmpl, err := template.ParseFiles("HTML/edit.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            log.Println("Error parsing template:", err)
            return
        }

        tmpl.Execute(w, map[string]string{
            "Username": username,
            "Email":    email,
            "Image":    image,
        })
    }
}
