package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    "regexp"
    "strings"
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
            media := r.FormValue("media")
            
            _, err := db.Exec("INSERT INTO post (user_id, text, media, date, category) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)", userID, text, media, category)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/myposts?status=success", http.StatusSeeOther)
            return
        }

        rows, err := db.Query(`
            SELECT p.post_id, p.text, p.media, p.date, p.category, u.username, u.image
            FROM post p
            JOIN user u ON p.user_id = u.id
            WHERE p.user_id = ?
            ORDER BY p.date DESC
        `, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var posts []struct {
            PostID   int
            Text     string
            Media    string
            MediaType string
            Date     string
            Category string
            Username string
            Image    string
        }
        
        for rows.Next() {
            var post struct {
                PostID   int
                Text     string
                Media    string
                MediaType string
                Date     string
                Category string
                Username string
                Image    string
            }
            err := rows.Scan(&post.PostID, &post.Text, &post.Media, &post.Date, &post.Category, &post.Username, &post.Image)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            post.MediaType = parseMediaType(post.Media)
            if post.MediaType == "youtube" {
                post.Media = embedYouTube(post.Media)
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
                MediaType string
                Date     string
                Category string
                Username string
                Image    string
            }
            StatusMessage string
        }{
            Posts:          posts,
            StatusMessage:  statusMessage,
        })
    }
}

// parseMediaType determines the type of media (image, youtube, or link)
func parseMediaType(media string) string {
    if media == "" {
        return "none"
    }
    if strings.Contains(media, "youtube.com") || strings.Contains(media, "youtu.be") {
        return "youtube"
    }
    if isImage(media) {
        return "image"
    }
    return "link"
}

// isImage checks if the URL is an image
func isImage(url string) bool {
    return strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif")
}

// embedYouTube converts a YouTube URL to an embeddable URL
func embedYouTube(url string) string {
    re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([\w-]+)`)
    match := re.FindStringSubmatch(url)
    if len(match) > 1 {
        return "https://www.youtube.com/embed/" + match[1]
    }
    return url
}

