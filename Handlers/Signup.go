package Forum

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
)

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

        
        //err
				if r.URL.Path != "/signup" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}



		
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			email := r.FormValue("email")
			image := r.FormValue("image") 
        



		    //err
			var errorMessage string

			if username == "" {
				errorMessage = "Username is required."
			} else if !isValidPassword(password) {
				errorMessage = "Password must be at least 8 characters long and include at least one letter, one number, and one special character."
			} else if email == "" {
				errorMessage = "Email is required."
			} else if !isValidEmail(email) {
				errorMessage = "Invalid email format."
			}


		    if errorMessage != "" {
							tmpl, err := template.ParseFiles("HTML/Signup.html")
							if err != nil {
								RenderErrorPage(w, http.StatusInternalServerError)
								return
							}
							tmpl.Execute(w, map[string]string{"ErrorMessage": errorMessage})
							return
			}



			var existingUser string
			err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&existingUser)
			if err != nil && err != sql.ErrNoRows {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}

            
			//err
			if existingUser != "" {
				tmpl, err := template.ParseFiles("HTML/Signup.html")
				if err != nil {
					RenderErrorPage(w, http.StatusInternalServerError)
					return
				}
				tmpl.Execute(w, map[string]string{"ErrorMessage": "Username already taken"})
				return
			}


			var existingEmail string
			err = db.QueryRow("SELECT email FROM user WHERE email = ?", email).Scan(&existingEmail)
			if err != nil && err != sql.ErrNoRows {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}

            
			//err
			if existingEmail != "" {
				tmpl, err := template.ParseFiles("HTML/Signup.html")
				if err != nil {
					RenderErrorPage(w, http.StatusInternalServerError)
					return
				}
				tmpl.Execute(w, map[string]string{"ErrorMessage": "Email already taken"})
				return
			}

			
			_, err = db.Exec("INSERT INTO user (username, password, email, image) VALUES (?, ?, ?, ?)", username, password, email, image)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError)
				return
			}

			
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} 

		tmpl, err := template.ParseFiles("HTML/Signup.html")
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		if isLetter(char) {
			hasLetter = true
		} else if isNumber(char) {
			hasNumber = true
		} else if isSpecial(char) {
			hasSpecial = true
		}
	}

	return hasLetter && hasNumber && hasSpecial
}

func isLetter(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func isSpecial(char rune) bool {
	return char == '@' || char == '$' || char == '!' || char == '%' || char == '*' || char == '?' || char == '&'
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}