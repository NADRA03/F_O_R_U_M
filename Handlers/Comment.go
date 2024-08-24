package Forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func CommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {


		//err
				if r.URL.Path != "/comment" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}


		session, _ := store.Get(r, "mysession")
        
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			username = defaultName
		}
	
		image, ok := session.Values["profileImage"].(string)
	
		if !ok || image == "" {
			image = defaultProfileImage 
		}

        postIDStr := r.URL.Query().Get("post_id")
        postID, err := strconv.Atoi(postIDStr)
        if err != nil {
            http.Error(w, "Invalid post ID", http.StatusBadRequest)
            RenderErrorPage(w, http.StatusBadRequest)
            return
        }
  
		if r.Method == http.MethodPost {
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			comment := r.FormValue("comment")


            //err
			if comment == "" {
				RenderErrorPage(w, http.StatusBadRequest) 
				return
            }


            userID, _ := session.Values["id"].(int)
			_, err := db.Exec("INSERT INTO comment (user_id, post_id, comment, date) VALUES (?, ?, ?, CURRENT_TIMESTAMP)", userID, postID, comment)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}
      
			http.Redirect(w, r, fmt.Sprintf("/comment?post_id=%d&status=success", postID), http.StatusSeeOther)
			return
		} 

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
		err = db.QueryRow(`
            SELECT p.post_id, p.text, p.media, p.date, p.category, u.username, u.image
            FROM post p
            JOIN user u ON p.user_id = u.id
            WHERE p.post_id = ?
        `, postID).Scan(&post.PostID, &post.Text, &post.Media, &post.Date, &post.Category, &post.Username, &post.Image)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				RenderErrorPage(w, http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				RenderErrorPage(w, http.StatusInternalServerError)
			}
			return
		}
		post.MediaType = parseMediaType(post.Media)
		if post.MediaType == "youtube" {
			post.Media = embedYouTube(post.Media)
		}

		// Retrieve comments including commenter details
		rows, err := db.Query(`
            SELECT c.comment_id, c.user_id, c.comment, c.date, u.username, u.image
            FROM comment c
            JOIN user u ON c.user_id = u.id
            WHERE c.post_id = ?
        `, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []struct {
			CommentID int
			UserID    int
			Comment   string
			Date      string
			Username  string
			Image     string
		}

		for rows.Next() {
			var comment struct {
				CommentID int
				UserID    int
				Comment   string
				Date      string
				Username  string
				Image     string
			}
			err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.Comment, &comment.Date, &comment.Username, &comment.Image)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}

        tmpl, err := template.ParseFiles("HTML/comment.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
      
        statusMessage := r.URL.Query().Get("status")
        tmpl.Execute(w, struct {
            Post           struct {
                PostID      int
                Text        string
                Media       string
                MediaType   string
                Date        string
                Category    string
                Username    string
                Image       string
            }
            Comments       []struct {
                CommentID    int
                UserID       int
                Comment      string
                Date         string
                Username     string
                Image        string
            }
            StatusMessage string
            Image         string
            Username      string
        }{  
            Username:       username,
            Image:          image,
            Post:           post,
            Comments:       comments,
            StatusMessage:  statusMessage,
        })
    }
}
