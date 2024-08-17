package Forum

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

// MyPostsHandler handles the "/myposts" route, allowing users to view and create posts.
func MyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session data
		session, _ := store.Get(r, "mysession")

		// Check if the user is authenticated; if not, redirect to the home page
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Get the user ID from the session
		userID, _ := session.Values["id"].(int)

		// Handle POST request: inserting a new post into the database
		if r.Method == http.MethodPost {
			text := r.FormValue("text")
			category := r.FormValue("category")
			media := r.FormValue("media")

			// Insert the new post into the 'post' table
			_, err := db.Exec("INSERT INTO post (user_id, text, media, date, category) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)", userID, text, media, category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				RenderErrorPage(w, http.StatusInternalServerError) 
				return
			}

			// Redirect the user back to the /myposts page with a success message
			http.Redirect(w, r, "/myposts?status=success", http.StatusSeeOther)
			return
		}

		// Fetch the user's posts from the database
		rows, err := db.Query(`
            SELECT p.post_id, p.text, p.media, p.date, p.category, u.username, u.image
            FROM post p
            JOIN user u ON p.user_id = u.id
            WHERE p.user_id = ?
            ORDER BY p.date DESC
        `, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		defer rows.Close() // Ensure rows are closed after query execution

		// Structure to hold post data
		var posts []struct {
			PostID    int
			Text      string
			Media     string
			MediaType string
			Date      string
			Category  string
			Username  string
			Image     string
		}

		// Iterate over the query results and populate the posts slice
		for rows.Next() {
			var post struct {
				PostID    int
				Text      string
				Media     string
				MediaType string
				Date      string
				Category  string
				Username  string
				Image     string
			}
			err := rows.Scan(&post.PostID, &post.Text, &post.Media, &post.Date, &post.Category, &post.Username, &post.Image)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				RenderErrorPage(w, http.StatusInternalServerError) 
				return
			}
			// Determine the media type (image, YouTube video, or link)
			post.MediaType = parseMediaType(post.Media)
			// Convert YouTube links to embeddable URLs
			if post.MediaType == "youtube" {
				post.Media = embedYouTube(post.Media)
			}
			// Add the post to the posts slice
			posts = append(posts, post)
		}

		// Parse the HTML template for displaying the posts
		tmpl, err := template.ParseFiles("HTML/myposts.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError)            
			return
		}

		// Retrieve any status message from the URL query parameters
		statusMessage := r.URL.Query().Get("status")
		// Render the template with the posts and the status message
		tmpl.Execute(w, struct {
			Posts []struct {
				PostID    int
				Text      string
				Media     string
				MediaType string
				Date      string
				Category  string
				Username  string
				Image     string
			}
			StatusMessage string
		}{
			Posts:         posts,
			StatusMessage: statusMessage,
		})

	}
}

// parseMediaType determines the type of media (image, YouTube, or link)
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
	// Regular expression to extract the video ID from a YouTube URL
	re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([\w-]+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return "https://www.youtube.com/embed/" + match[1]
	}
	return url
}
