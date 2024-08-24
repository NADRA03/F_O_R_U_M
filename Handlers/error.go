package Forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderErrorPage(w http.ResponseWriter, statusCode int) {
	errorPage := fmt.Sprintf("HTML/%d.html", statusCode)
	tmpl, err := template.ParseFiles(errorPage)
	if err != nil {
		http.Error(w, "Error loading error template", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	tmpl.Execute(w, nil)
}
