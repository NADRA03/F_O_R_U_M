package Forum

import (
    "database/sql"
    "net/http"
)

func AddPostHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        

        //err
        if r.URL.Path != "/addpost" {
            RenderErrorPage(w, http.StatusNotFound)
            return
        }

        session, _ := store.Get(r, "mysession")

        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        userID, _ := session.Values["id"].(int)

        if r.Method == http.MethodPost {
            text := r.FormValue("text")
            category := r.FormValue("category")

            
            //err
            if text == "" || category == "" {
                            RenderErrorPage(w, http.StatusBadRequest) 
                            return
            }

            _, err := db.Exec("INSERT INTO post (user_id, text, media, date, category) VALUES (?, ?, '', CURRENT_TIMESTAMP, ?)", userID, text, category)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                RenderErrorPage(w, http.StatusInternalServerError)  
                return
            }

            http.Redirect(w, r, "/?status=success", http.StatusSeeOther)
            return
        }
        RenderErrorPage(w, http.StatusMethodNotAllowed)
    }
}
