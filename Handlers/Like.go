package Forum

import (
    "database/sql"
    "net/http"
)

func LikeHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")

        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        userID, _ := session.Values["id"].(int)
        postID := r.URL.Query().Get("post_id")
        
		var existingLikeCount int
        err := db.QueryRow("SELECT COUNT(*) FROM `like` WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingLikeCount)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }

        if existingLikeCount > 0 {
            // User has already liked this post
            http.Redirect(w, r, "/myposts", http.StatusSeeOther)
            return
        }
        _, err = db.Exec("INSERT INTO `like` (user_id, post_id) VALUES (?, ?)", userID, postID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
        
        http.Redirect(w, r, "/myposts?StatusLike=success", http.StatusSeeOther)
    }
}
