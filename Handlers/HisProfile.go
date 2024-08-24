package Forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type Post struct {
	PostID    int
	Username  string
	MediaType string
	Image     string
	Text      string
	Media     string
	Date      string
	Category  string
	User_id int
}

type ProfilePageData struct {
	Username string
	Posts    []Post
	Name string
	ProfileImage string
	Followings    int
	Followers     int
	ProfileOwnerImage string
}

func ViewProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        

		//err
			if r.URL.Path != "/profile" {
				RenderErrorPage(w, http.StatusNotFound)
				return
			}


		session, _ := store.Get(r, "mysession")
		username := r.URL.Query().Get("username")
  



        //err
		if username == "" {
			RenderErrorPage(w, http.StatusBadRequest)
			return
		}




		name, ok := session.Values["username"].(string)
		if !ok || name == "" {
			name = defaultName
		}
		profileImage, ok := session.Values["profileImage"].(string)
		if !ok || profileImage == "" {
			profileImage = defaultProfileImage 
		}
        

        var profileOwnerImage string
		err := db.QueryRow("SELECT image FROM user WHERE username = ?", username).Scan( &profileOwnerImage)
		//err
		if err == sql.ErrNoRows {
			RenderErrorPage(w, http.StatusNotFound) 
			return
		} else if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}


		var followerCount int
		var userID int
        err = db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
        err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE user_id = ?", userID).Scan(&followerCount)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

		var followingCount int
        err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", userID).Scan(&followingCount)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }
		query := `
            SELECT post.post_id, user.username, user.id, user.image, post.text, post.media, post.date, post.category
            FROM post
            JOIN user ON post.user_id = user.id
            WHERE user.username = ?
			ORDER BY post.date DESC
        `
		rows, err := db.Query(query, username)
		if err != nil {
			log.Printf("Error querying posts for username %s: %v", username, err)
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.PostID, &post.Username, &post.User_id, &post.Image, &post.Text, &post.Media, &post.Date, &post.Category); err != nil {
				log.Printf("Error scanning post: %v", err)
				RenderErrorPage(w, http.StatusInternalServerError)
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
            RenderErrorPage(w, http.StatusInternalServerError) 
			return
		}

		tmpl, err := template.ParseFiles("HTML/profile.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
            RenderErrorPage(w, http.StatusInternalServerError) 
			return
		}

		data := ProfilePageData{
			Username: username,
			Posts:    posts,
			Name:     name, 
            ProfileImage: profileImage,
			Followings: followingCount,
			Followers: followerCount,
			ProfileOwnerImage: profileOwnerImage, 
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error executing template: %v", err)
            RenderErrorPage(w, http.StatusInternalServerError) 
		}
	}
}
