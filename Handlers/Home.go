package Forum

import (
    "html/template"
    "net/http"
    _ "github.com/gorilla/sessions"
)

const defaultProfileImage = "https://www.strasys.uk/wp-content/uploads/2022/02/Depositphotos_484354208_S.jpg" // Default image link

func RootHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "mysession")

     // Check if the user is authenticated
     if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    name, _ := session.Values["username"].(string)
    profileImage, ok := session.Values["profileImage"].(string)

    if !ok || profileImage == "" {
        profileImage = defaultProfileImage // Use default image if not set
    }

    tmpl, err := template.ParseFiles("HTML/Home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        RenderErrorPage(w, http.StatusInternalServerError) 
        
        return
    }
    statusMessage := r.URL.Query().Get("status")
    tmpl.Execute(w, struct {
        Name         string
        ProfileImage string
        StatusMessage string
    }{
        Name:         name,
        ProfileImage: profileImage,
        StatusMessage: statusMessage,
    })
}


