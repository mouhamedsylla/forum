package server

import (
	"forum/Api/controllers"
	"forum/Api/services"
	"forum/server/middleware"
	"forum/server/router"
	"log"
	"net/http"
)

type Server struct {
	Router *router.Router
}

func NewServer() *Server {
	return &Server{
		Router: router.NewRouter(),
	}
}

func (s *Server) SessionDataInit() {
	services.Global_SessionManager.InitSessionMap()
}

func (s *Server) ConfigureRoutes() {
	controllers := controllers.NewControllers()
	controllers.PostInfos()
	s.Router.SetDirectory("/assets/", "./App/internal/assets/")
	s.Router.Method(http.MethodGet).Handler("/assets/", s.Router.StaticServe())
	s.Router.Method(http.MethodGet).Handler("/", controllers.Forum())
	s.Router.Method(http.MethodPost).Handler("/login", controllers.Login())
	s.Router.Method(http.MethodPost).Handler("/post", controllers.CreatePost())
	s.Router.Method(http.MethodGet).Middleware(middleware.SessionMiddleware).Handler("/home", controllers.Home())
	s.Router.Method(http.MethodGet).Handler("/connect", controllers.Connection())
	s.Router.Method(http.MethodGet).Middleware(middleware.SessionMiddleware).Handler("/session", controllers.SessionUser())
	s.Router.Method(http.MethodPost).Handler("/register", controllers.Register())
	s.Router.Method(http.MethodPost).Handler("/category", controllers.AddCategories())
	s.Router.Method(http.MethodGet).Handler("/getposts", controllers.GetPosts())
	s.Router.Method(http.MethodPost).Middleware(middleware.SessionMiddleware).Handler("/reaction", controllers.Reaction())
	s.Router.Method(http.MethodPost).Middleware(middleware.SessionMiddleware).Handler("/logout", controllers.Logout())
	s.Router.Method(http.MethodPost).Middleware(middleware.SessionMiddleware).Handler("/comment", controllers.Comment())
	s.Router.Method(http.MethodPost).Handler("/getcomment", controllers.GetComments())
	s.Router.Method(http.MethodGet).Handler("/error", controllers.RenderErrorPage())
	s.Router.Method(http.MethodPost).Handler("/getcommentuser", controllers.CommentUser())
	s.Router.Method(http.MethodPost).Middleware(middleware.SessionMiddleware).Handler("/commentreaction", controllers.CommentReaction())
	s.Router.Method(http.MethodGet).Middleware(middleware.SessionMiddleware).Handler("/getlikedpost", controllers.GetLikedPost())
}

func (s *Server) StartServer(port string) error {
	s.SessionDataInit()
	s.ConfigureRoutes()
	log.Printf("http://localhost:%v\n", port)
	return http.ListenAndServe(":"+port, s.Router)
}
