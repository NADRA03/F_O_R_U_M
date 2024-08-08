package Forum

import (
    "html/template"
    "net/http"
    _ "github.com/gorilla/sessions"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "mysession")

    // Check if the user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    name, _ := session.Values["username"].(string)

    tmpl, err := template.ParseFiles("HTML/Home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    statusMessage := r.URL.Query().Get("status")
    tmpl.Execute(w, struct {
        Name string
        StatusMessage string
    }{
        Name: name,
        StatusMessage:  statusMessage,
    })
}


