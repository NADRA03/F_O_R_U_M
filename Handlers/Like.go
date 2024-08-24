package Forum

import (
    "database/sql"
    "net/http"
    "strconv"
)

func LikeHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "mysession")
        

        //err
				if r.URL.Path != "/like" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}


        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        userID, _ := session.Values["id"].(int)
        postID := r.URL.Query().Get("post_id")


        //err
                if postID == "" {
                    RenderErrorPage(w, http.StatusBadRequest) 
                    return
                }
                
        //err
                postIDint, err := strconv.Atoi(postID)
                if err != nil || postIDint <= 0 {
                    RenderErrorPage(w, http.StatusBadRequest) 
                    return
                }
        

        //err        
        var postExists bool
        err = db.QueryRow("SELECT COUNT(*) > 0 FROM post WHERE post_id = ?", postIDint).Scan(&postExists)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }
        if !postExists {
            RenderErrorPage(w, http.StatusNotFound)
            return
        }
        
    
        
		var existingLikeCount int
        err = db.QueryRow("SELECT COUNT(*) FROM `like` WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingLikeCount)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
        

        //fix
        if existingLikeCount > 0 {
            _, err = db.Exec("DELETE FROM `like` WHERE user_id = ? AND post_id = ?", userID, postIDint)
            if err != nil {
                RenderErrorPage(w, http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Like removed"))
            return
        }


        _, err = db.Exec("INSERT INTO `like` (user_id, post_id) VALUES (?, ?)", userID, postID)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }


        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Like added"))
    }
}
