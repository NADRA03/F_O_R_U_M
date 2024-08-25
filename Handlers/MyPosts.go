package Forum

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"strconv"
)


func MyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {


        //err
				if r.URL.Path != "/myposts" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}


		
		session, _ := store.Get(r, "mysession")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
 
		
		userID, _ := session.Values["id"].(int)
		username, _ := session.Values["username"].(string)
		image, _ := session.Values["profileImage"].(string)
		var followerCount int
        err := db.QueryRow("SELECT COUNT(*) FROM followers WHERE user_id = ?", userID).Scan(&followerCount)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

		var followingCount int
        err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", userID).Scan(&followingCount)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }


		action := r.URL.Query().Get("action")
		if action == "delete" {
			postIDStr := r.URL.Query().Get("post_id")
			if postIDStr == "" {
				RenderErrorPage(w, http.StatusBadRequest)
				return
			}
		
			postID, err := strconv.Atoi(postIDStr)
			if err != nil {
				RenderErrorPage(w, http.StatusBadRequest) 
				return
			}
		
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE post_id = ? AND user_id = ?)", postID, userID).Scan(&exists)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}
		
			if !exists {
				RenderErrorPage(w, http.StatusNotFound)
				return
			}
		
			_, err = db.Exec("DELETE FROM post WHERE post_id = ? AND user_id = ?", postID, userID)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}
		
			// Redirect after deletion
			http.Redirect(w, r, "/myposts?status=deleted", http.StatusSeeOther)
			return
		}



		if r.Method == http.MethodPost {
			text := r.FormValue("text")
			category := r.FormValue("category")
			media := r.FormValue("media")
			
			_, err := db.Exec("INSERT INTO post (user_id, text, media, date, category) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)", userID, text, media, category)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError) 
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
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		defer rows.Close() 


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
				RenderErrorPage(w, http.StatusInternalServerError) 
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
            RenderErrorPage(w, http.StatusInternalServerError)            
			return
		}


		statusMessage := r.URL.Query().Get("status")
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
			Followings    int
			Followers     int
			Image         string
			Username      string 
		}{   
			Followings: followingCount,
			Followers: followerCount,
			Posts:         posts,
			StatusMessage: statusMessage,
			Image:         image,
			Username:      username,
		})

	}
}


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


func isImage(url string) bool {
	return strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif")
}


func embedYouTube(url string) string {
	re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/shorts/)([\w-]+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return "https://www.youtube.com/embed/" + match[1]
	}
	return url
}
