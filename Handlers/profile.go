package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    // "github.com/gorilla/sessions"
)

// var store = sessions.NewCookieStore([]byte("super-secret-key"))

type ProfileData struct {
    Username string
    Email    string
    Image    string
}

func ProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")
        auth, ok := session.Values["authenticated"].(bool)
        
        if !ok || !auth {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        username, _ := session.Values["username"].(string)
        
        var profileData ProfileData
        err := db.QueryRow("SELECT username, email, image FROM user WHERE username = ?", username).
            Scan(&profileData.Username, &profileData.Email, &profileData.Image)
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        } else if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
            // http.Redirect(w, r, "/signup", http.StatusSeeOther)
            return
        }

        tmpl, err := template.ParseFiles("HTML/profile.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        tmpl.Execute(w, profileData)
    }
}
