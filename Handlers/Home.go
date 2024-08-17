package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    _ "github.com/gorilla/sessions"

)

var defaultProfileImage = "https://www.strasys.uk/wp-content/uploads/2022/02/Depositphotos_484354208_S.jpg" 
var defaultName = "anonymous"

func RootHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "mysession")
    
    name, ok := session.Values["username"].(string)
    if !ok || name == "" {
        name = defaultName
    }

    profileImage, ok := session.Values["profileImage"].(string)

    if !ok || profileImage == "" {
        profileImage = defaultProfileImage 
    }
    
    category := r.URL.Query().Get("category")

    query := `
        SELECT p.post_id, p.text, p.media, p.date, p.category, u.username, u.image
        FROM post p
        JOIN user u ON p.user_id = u.id
    `
    if category != "" {
        query += `WHERE p.category = ? `
    }
    query += `ORDER BY p.date DESC`

    var rows *sql.Rows
    var err error
    if category != "" {
        rows, err = db.Query(query, category)
    } else {
        rows, err = db.Query(query)
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
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        post.MediaType = parseMediaType(post.Media)
        if post.MediaType == "youtube" {
            post.Media = embedYouTube(post.Media)
        }
        posts = append(posts, post)
    }
    tmpl, err := template.ParseFiles("HTML/Home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
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


