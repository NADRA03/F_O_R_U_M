package Forum

import (
    "html/template"
    "net/http"
    _ "github.com/gorilla/sessions"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "mysession")
    if r.URL.Path == "/" {
        http.Redirect(w, r, "/home", http.StatusFound)
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


