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

        
        //err
				if r.URL.Path != "/login" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}


        session, _ := store.Get(r, "mysession")
        
        if r.URL.Query().Get("logout") == "true" {
            session.Options.MaxAge = -1
            session.Save(r, w)

            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        
        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            password := r.FormValue("password")

            var id int
            var storedPassword string
            var image string
            err := db.QueryRow("SELECT id, password, image FROM user WHERE username = ?", username).Scan(&id, &storedPassword, &image)
            if err != nil || password != storedPassword {
                http.Redirect(w, r, "/login?error=invalid", http.StatusSeeOther)
                return
            }

            session.Values["id"] = id
            session.Values["authenticated"] = true
            session.Values["username"] = username
            session.Values["profileImage"] = image
            session.Save(r, w)

            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        } 

        tmpl, err := template.ParseFiles("HTML/Login.html")
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
        tmpl.Execute(w, nil)
    }
}


