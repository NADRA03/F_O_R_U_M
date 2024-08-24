package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    // "fmt"

)

var moderators = map[string]bool{
    "moderator": true, // Add your moderator usernames here
}

func ModeratorHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        
        //err
				if r.URL.Path != "/moderator" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}



        session, _ := store.Get(r, "mysession")

        // Check if the user is authenticated
        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        username, _ := session.Values["username"].(string)
        if !moderators[username] {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        // if r.Method == http.MethodPost {
        //     err := createTables(db)
        //     if err != nil {
        //         http.Error(w, err.Error(), http.StatusInternalServerError)
        //         return
        //     }
        //     fmt.Fprintln(w, "Tables created successfully.")
        //     return
        // }

        tables, err := fetchTables(db)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }

        tmpl, err := template.ParseFiles("HTML/Moderator.html")
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
        tmpl.Execute(w, struct {
            Tables []Table
        }{Tables: tables})
    }
}
