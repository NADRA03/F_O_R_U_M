package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
)

func MyPostsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")

        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
      
        userID, _ := session.Values["id"].(int)

        if r.Method == http.MethodPost {
            text := r.FormValue("text")
            category := r.FormValue("category")

            _, err := db.Exec("INSERT INTO post (user_id, text, media, date, category) VALUES (?, ?, '', CURRENT_TIMESTAMP, ?)", userID, text, category)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/myposts?status=success", http.StatusSeeOther)
            return
        }

        rows, err := db.Query("SELECT post_id, text, media, date, category FROM post WHERE user_id = ?", userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var posts []struct {
            PostID   int
            Text     string
            Media    string
            Date     string
            Category string
        }

        for rows.Next() {
            var post struct {
                PostID   int
                Text     string
                Media    string
                Date     string
                Category string
            }
            err := rows.Scan(&post.PostID, &post.Text, &post.Media, &post.Date, &post.Category)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            posts = append(posts, post)
        }

        tmpl, err := template.ParseFiles("HTML/myposts.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        statusMessage := r.URL.Query().Get("status")
        tmpl.Execute(w, struct {
            Posts          []struct {
                PostID   int
                Text     string
                Media    string
                Date     string
                Category string
            }
            StatusMessage string
        }{
            Posts:          posts,
            StatusMessage:  statusMessage,
        })
    }
}

