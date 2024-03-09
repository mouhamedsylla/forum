package controllers

import (
	"forum/Api/models"
	"forum/Api/services"
	"forum/utils"
	"html/template"
	"net/http"
)

type AllPosts struct {
	Posts []*models.PostUI
}

func (c *Controllers) Forum() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message

		tmpl, err := template.ParseFiles("../App/internal/assets/index.html")
		if err != nil {
			message.Error = err.Error()
			DataError.Message = "page not found"
			DataError.ErrorCode = http.StatusNotFound
			http.Redirect(w, r, "http://localhost:8080/error", http.StatusSeeOther)
			return
		}

		sessionID, err := r.Cookie("forum")
		if err == nil {
			sessionID := services.Global_SessionManager.GetSession(sessionID.Value)
			if sessionID != nil {
				http.Redirect(w, r, "http://localhost:8080/home", http.StatusSeeOther)
				return
			}
		}
		posts := AllPosts{
			Posts: PostsUI,
		}
		tmpl.Execute(w, posts)
	})
}

func (c *Controllers) Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		tmpl, err := template.ParseFiles("../App/internal/assets/home.html")

		if err != nil {
			message.Error = err.Error()
			DataError.Message = "page not found"
			DataError.ErrorCode = http.StatusNotFound
			http.Redirect(w, r, "http://localhost:8080/error", http.StatusSeeOther)
			return
		}

		posts := AllPosts{
			Posts: PostsUI,
		}
		tmpl.Execute(w, posts)
	})
}

func (c *Controllers) SessionUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("user").(*models.Session)
		utils.RespondWithJSON(w, session, http.StatusOK)
	})
}
