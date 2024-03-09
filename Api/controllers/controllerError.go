package controllers

import (
	"forum/Api/models"
	"html/template"
	"net/http"
)

var (
	DataError = models.ErrorData {}
)

func (c *Controllers) RenderErrorPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("../App/internal/assets/error.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(DataError.ErrorCode)
		err = tmpl.Execute(w, DataError)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

}
