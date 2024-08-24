package Forum

import (
    "database/sql"
    "net/http"
    "strconv"
)

func FollowHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

         
        //err
				if r.URL.Path != "/follow" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}




        session, _ := store.Get(r, "mysession")
        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        userID := r.URL.Query().Get("user_id")
        followerID, _ := session.Values["id"].(int)
        

        //err
        if userID == "" {
			RenderErrorPage(w, http.StatusBadRequest) 
			return
		}
        
        //err
        userIDint, err := strconv.Atoi(userID)
		if err != nil || userIDint <= 0 {
			RenderErrorPage(w, http.StatusBadRequest) 
			return
		}
        
        
        //err
        var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)", userID).Scan(&exists)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		if !exists {
			RenderErrorPage(w, http.StatusNotFound) 
			return
		}
        
        
        var existingFollowCount int
        err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?", userID, followerID).Scan(&existingFollowCount)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }
        


        //fix
        if existingFollowCount > 0 {
            _, err = db.Exec("DELETE FROM followers WHERE user_id = ? AND follower_id = ?", userIDint, followerID)
            if err != nil {
                RenderErrorPage(w, http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Follow removed"))
            return
        }



        _, err = db.Exec("INSERT INTO followers (user_id, follower_id) VALUES (?, ?)", userID, followerID)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Follow added"))
    }
}
