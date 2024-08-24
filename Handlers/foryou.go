package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    _ "github.com/gorilla/sessions"

)

func ForyouHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {


    //err
            if r.URL.Path != "/foryou" {
                RenderErrorPage(w, http.StatusNotFound)
                return
    }



    session, _ := store.Get(r, "mysession")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
    
	userID,_ := session.Values["id"].(int)
    category := r.URL.Query().Get("category")

	query := `
	SELECT p.post_id, p.text, p.media, p.date, p.category, u.username, u.image
	FROM post p
	JOIN user u ON p.user_id = u.id
	WHERE p.user_id IN (
		SELECT f.user_id
		FROM followers f
		WHERE f.follower_id = ?
	)
    `
	if category != "" {
	query += ` AND p.category = ? `
	}
	query += `ORDER BY p.date DESC`

	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = db.Query(query, userID, category)
	} else {
		rows, err = db.Query(query, userID)
	}
    if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
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
    tmpl, err := template.ParseFiles("HTML/foryou.html")
    if err != nil {
        RenderErrorPage(w, http.StatusInternalServerError) 
        return
    }
    statusMessage := r.URL.Query().Get("status")
    tmpl.Execute(w, struct {
        Posts          []struct {
            PostID    int
            Text      string
            Media     string
            MediaType string
            Date      string
            Category  string
            Username  string
            Image     string

        }
        Name         string
        ProfileImage string
        StatusMessage string
    }{  
        Posts:         posts,
        Name:          name,
        ProfileImage: profileImage,
        StatusMessage: statusMessage,

    })
}
}
