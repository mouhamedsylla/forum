package controllers

import (
	"forum/Api/models"
	"forum/Api/services"
	"forum/utils"
	"net/http"
	"regexp"
	"strconv"
)

// This `Register()` function is a method defined on the `Controllers` struct. It is used to handle the
// registration process for a user in a web application. Here's a breakdown of what the function does:
func (handler *Controllers) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.User{})

		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}

		user := result.(*models.User)
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		validEmail := regexp.MustCompile(emailRegex).MatchString(user.Email)

		if !validEmail || user.Username == "" {
			message.Error = "Invalid email format"
			utils.RespondWithJSON(w, message, http.StatusBadRequest)
			return
		}
		models.CryptPassword(user)
		err = handler.Storage.Insert(*user)

		if err != nil {
			message.Error = "email already used"
			utils.RespondWithJSON(w, message, http.StatusInternalServerError)
			return
		}
		message.Message = "registration successful"
		utils.RespondWithJSON(w, message, http.StatusOK)

	})
}

// The `func (handler *Controllers) Login() http.Handler` function is a method defined on the
// `Controllers` struct. It is used to handle the login process for a user in a web application. Here's
// a breakdown of what the function does:
func (handler *Controllers) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.User{})

		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}

		toAuthenticate := *result.(*models.User)
		rslt := handler.Storage.Select(toAuthenticate, "Email", toAuthenticate.Email).([]models.User)
		if len(rslt) == 0 {
			message.Error = "this user does't exist"
			utils.RespondWithJSON(w, message, http.StatusNotFound)
			return
		}

		user := rslt[0]
		if err = handler.Auth.Authenticate(user.Password, &toAuthenticate); err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, http.StatusUnauthorized)
			return
		}

		handler.Storage.DeleteExipireSession(user.Id)
		services.Global_SessionManager.DeleteSessionWithUserID(user.Id)
		services.Global_SessionManager.GenerateSession(&w, &user)
		reactions := handler.Storage.Select(models.ReactionPost{}, "User_id", user.Id).([]models.ReactionPost)

		if len(reactions) != 0 {
			for _, r := range reactions {
				key := strconv.Itoa(r.PostId) + "-" + strconv.Itoa(user.Id)
				Post_reaction[key] = r.Value
			}
		}
		message.Message = "you are now logged"
		utils.RespondWithJSON(w, message, http.StatusOK)
	})
}
