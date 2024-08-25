package Forum

import (
	"database/sql"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	// "github.com/gorilla/sessions"
)

func EditProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {


        //err
				if r.URL.Path != "/settings" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}



        session, _ := store.Get(r, "mysession")
        userID, ok := session.Values["id"].(int)
		// fmt.Println("userID: ", userID)
        if !ok || session.Values["authenticated"] != true {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            email := r.FormValue("email")
            password := r.FormValue("password")
            image := r.FormValue("image")


            //err
            var existingID int
            err := db.QueryRow("SELECT id FROM user WHERE id = ?", userID).Scan(&existingID)
            if err != nil {
                if err == sql.ErrNoRows {
                    RenderErrorPage(w, http.StatusNotFound)
                } else {
                    RenderErrorPage(w, http.StatusInternalServerError)
                }
                log.Println("User ID not found or other error:", userID, err)
                return
            }
            
            //err
            if username == "" || email == "" || password == "" {
				RenderErrorPage(w, http.StatusBadRequest)
				return
			}

            		
		var existingUser string
		err = db.QueryRow("SELECT username FROM user WHERE username = ? AND id != ?", username, userID).Scan(&existingUser)
		if err != nil && err != sql.ErrNoRows {
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		if existingUser != "" {
			RenderErrorPage(w, http.StatusBadRequest)
			return
		}

		
		var existingEmail string
		err = db.QueryRow("SELECT email FROM user WHERE email = ? AND id != ?", email, userID).Scan(&existingEmail)
		if err != nil && err != sql.ErrNoRows {
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		if existingEmail != "" {
			RenderErrorPage(w, http.StatusBadRequest)
			return
		}

            _, err = db.Exec("UPDATE user SET username = ?, email = ?, password = ?, image = ? WHERE id = ?",
                username, email, password, image, userID)
            if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError)
                log.Println("Error updating user:", err)
                return
            }

            session.Values["username"] = username
            session.Values["profileImage"] = image
            session.Save(r, w)
            http.Redirect(w, r, "/settings", http.StatusSeeOther)
            return
        } 

        var username, email, image string
        err := db.QueryRow("SELECT username, email, image FROM user WHERE id = ?", userID).Scan(&username, &email, &image)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            log.Println("Error retrieving user information:", err)
            return
        }

        tmpl, err := template.ParseFiles("HTML/edit.html")
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            log.Println("Error parsing template:", err)
            return
        }

        tmpl.Execute(w, map[string]string{
            "Username": username,
            "Email":    email,
            "Image":    image,
        })
    }
}
