package controllers

import (
	"forum/Api/models"
	"forum/Api/services"
	"forum/utils"
	"html/template"
	"net/http"
)

func (c *Controllers) Connection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		sessionID, err := r.Cookie("forum")
		if err == nil {
			sessionID := services.Global_SessionManager.GetSession(sessionID.Value)
			if sessionID != nil {
				http.Redirect(w, r, "http://localhost:8080/home", http.StatusSeeOther)
				return
			}
		}
		tmpl, err := template.ParseFiles("../App/internal/assets/connection.html")

		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, http.StatusNotFound)
			return
		}
		tmpl.Execute(w, nil)
	})
}

func (c *Controllers) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.Session{})
		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}

		session := result.(*models.Session)
		if session != nil {
			c.Storage.DeleteExipireSession(session.User.Id)
			services.Global_SessionManager.DeleteSessionWithUserID(session.User.Id)
		}
		message.Message = "logout done"
		utils.RespondWithJSON(w, message, http.StatusOK)
	})
}
