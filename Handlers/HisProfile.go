package Forum

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
)

type Post struct {
    PostID       int
    Username     string
	MediaType    string
    Image        string
    Text         string
    Media        string
    Date         string
    Category     string
}

type ProfilePageData struct {
    Username string
    Posts    []Post
}

func ViewProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username := r.URL.Query().Get("username")
        if username == "" {
            http.Error(w, "Username is required", http.StatusBadRequest)
            return
        }

        query := `
            SELECT post.post_id, user.username, user.image, post.text, post.media, post.date, post.category
            FROM post
            JOIN user ON post.user_id = user.id
            WHERE user.username = ?
			ORDER BY post.date DESC
        `
        rows, err := db.Query(query, username)
        if err != nil {
            log.Printf("Error querying posts for username %s: %v", username, err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var posts []Post
        for rows.Next() {
            var post Post
            if err := rows.Scan(&post.PostID, &post.Username, &post.Image, &post.Text, &post.Media, &post.Date, &post.Category); err != nil {
                log.Printf("Error scanning post: %v", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
			post.MediaType = parseMediaType(post.Media)
            if post.MediaType == "youtube" {
                post.Media = embedYouTube(post.Media)
            }
            posts = append(posts, post)
        }

        if err := rows.Err(); err != nil {
            log.Printf("Error iterating over posts: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        tmpl, err := template.ParseFiles("HTML/profile.html")
        if err != nil {
            log.Printf("Error parsing template: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        data := ProfilePageData{
            Username: username,
            Posts:    posts,
        }

        if err := tmpl.Execute(w, data); err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }
}


