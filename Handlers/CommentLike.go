package Forum

import (
    "database/sql"
    "net/http"
    "strconv"
)

func CommentLikeHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")

        // Check the URL path
        if r.URL.Path != "/comment-like" {
            RenderErrorPage(w, http.StatusNotFound)
            return
        }

        // Ensure the user is authenticated
        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        userID, _ := session.Values["id"].(int)
        commentID := r.URL.Query().Get("comment_id")

        // Validate comment ID
        if commentID == "" {
            RenderErrorPage(w, http.StatusBadRequest)
            return
        }

        commentIDint, err := strconv.Atoi(commentID)
        if err != nil || commentIDint <= 0 {
            RenderErrorPage(w, http.StatusBadRequest)
            return
        }

        // Check if the comment exists
        var commentExists bool
        err = db.QueryRow("SELECT COUNT(*) > 0 FROM comment WHERE comment_id = ?", commentIDint).Scan(&commentExists)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }
        if !commentExists {
            RenderErrorPage(w, http.StatusNotFound)
            return
        }

        // Check if the user already liked this comment
        var existingLikeCount int
        err = db.QueryRow("SELECT COUNT(*) FROM comment_like WHERE user_id = ? AND comment_id = ?", userID, commentIDint).Scan(&existingLikeCount)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

        // If the like exists, remove it
        if existingLikeCount > 0 {
            _, err = db.Exec("DELETE FROM comment_like WHERE user_id = ? AND comment_id = ?", userID, commentIDint)
            if err != nil {
                RenderErrorPage(w, http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Like removed"))
            return
        }

        // Otherwise, add the like
        _, err = db.Exec("INSERT INTO comment_like (user_id, comment_id) VALUES (?, ?)", userID, commentIDint)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Like added"))
    }
}
