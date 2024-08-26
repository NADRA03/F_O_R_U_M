// main.go
package main

import (
    "database/sql"
    "fmt"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "Forum/Handlers"
)

func main() {
    db, err := sql.Open("sqlite3", "./Forum.db")
    if err != nil {
        fmt.Println("Failed to open database:", err)
        return
    }
    defer db.Close()

    http.HandleFunc("/", Forum.RootHandler(db))
    http.HandleFunc("/moderator", Forum.ModeratorHandler(db))
    http.HandleFunc("/login", Forum.LoginHandler(db))
    http.HandleFunc("/signup", Forum.SignupHandler(db))
    http.HandleFunc("/myposts", Forum.MyPostsHandler(db))
    http.HandleFunc("/mylikes", Forum.MyLikesHandler(db))
    http.HandleFunc("/like", Forum.LikeHandler(db))
    http.HandleFunc("/comment", Forum.CommentHandler(db))
    http.HandleFunc("/addpost",Forum.AddPostHandler(db))
    http.HandleFunc("/settings", Forum.EditProfileHandler(db))  
    http.HandleFunc("/profile",Forum.ViewProfileHandler(db))
    http.HandleFunc("/follow", Forum.FollowHandler(db))
    http.HandleFunc("/foryou", Forum.ForyouHandler(db))
    http.HandleFunc("/comment-like", Forum.CommentLikeHandler(db))
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
