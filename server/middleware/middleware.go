package middleware

import (
	"context"
	"forum/Api/controllers"
	"forum/Api/services"
	"net/http"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("forum")
		if err != nil {

			controllers.DataError.Message = "Unauthorized"
			controllers.DataError.ErrorCode = http.StatusUnauthorized

			http.Redirect(w, r,"http://localhost:8080/error", http.StatusSeeOther)
			return
		}
		session := services.Global_SessionManager.GetSession(sessionID.Value)
		if session == nil {
			http.Redirect(w, r, "http://localhost:8080/connect", http.StatusSeeOther)
			return
		}

		user := services.Global_SessionManager.Sessions[sessionID.Value]

		// Ajoutez le user à la requête pour l'utiliser dans les gestionnaires
		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next.ServeHTTP(w, r)
	})
}
