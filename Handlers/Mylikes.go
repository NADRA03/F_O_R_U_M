package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
)

func MyLikesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")

        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        userID, _ := session.Values["id"].(int)
        username, _ := session.Values["username"].(string)
        image, _ := session.Values["profileImage"].(string)
        // Query to select posts liked by the user
        rows, err := db.Query(`
            SELECT p.post_id, p.text, p.media, p.date, p.category 
            FROM post p
            JOIN like l ON p.post_id = l.post_id
            WHERE l.user_id = ?
        `, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError) 
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
                RenderErrorPage(w, http.StatusInternalServerError) 
                return
            }
            posts = append(posts, post)
        }

        tmpl, err := template.ParseFiles("HTML/mylikes.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }

        tmpl.Execute(w, struct {
            Posts          []struct {
                PostID   int
                Text     string
                Media    string
                Date     string
                Category string
            }
            Username      string
            Image         string
        }{
            Posts:          posts,
            Username:       username,
            Image:          image,
        })
    }
}
