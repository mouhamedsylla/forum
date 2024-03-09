package boot

import (
	"fmt"
	"forum/Api/models"
	"forum/App/internal/components"
	directive "forum/App/internal/directives"
	"forum/App/internal/utils"
	"net/http"
	"syscall/js"
)

// The type `App` contains a slice of `directive.Directive` and a pointer to `models.Session`.
// @property {[]directive.Directive} Directives - The `Directives` property in the `App` struct is a
// slice of `directive.Directive` objects. This property likely holds a collection of directives that
// provide instructions or guidelines for the application to follow.
// @property Session - The `Session` property in the `App` struct is a pointer to an instance of the
// `Session` struct from the `models` package. This allows the `App` struct to hold a reference to a
// session object that contains relevant data and functionality for the application.
type App struct {
	Directives []directive.Directive
	Session    *models.Session
}

// The NewApp function in Go returns a new instance of the App struct.
func NewApp() *App {
	return &App{}
}

// The `BootstrapApplication` method in the `App` struct is iterating over each directive stored in the
// `Directives` slice. For each directive, it retrieves the corresponding elements from the DOM using
// the selector defined in the directive. It then sets the element for the directive and initializes it
// by calling the `Init` method on the directive. This process allows the application to apply the
// necessary behavior or functionality defined by each directive to the selected elements on the
// webpage.
func (app *App) BootstrapApplication() {
	fmt.Println("Init App")
	for _, directive := range app.Directives {
		elements := js.Global().Get("document").Call("querySelectorAll", directive.Selector())
		for i := 0; i < elements.Length(); i++ {
			directive.SetElement(elements.Index(i))
			directive.Init()
		}
	}
}

// The `func (app *App) InitSession() error` function in the provided Go code snippet is responsible
// for initializing the session data for the application. Here's a breakdown of what this function
// does:
func (app *App) InitSession() error {
	fmt.Println("ok Session")
	errChan := make(chan error)
	utils.GetData("session", func(r *http.Response, err error) {
		//defer fmt.Println(err)
		if err = utils.DecodeResponse(r, &app.Session); err != nil {
			errChan <- err
		}
		close(errChan)
	})
	return <-errChan
}

// The `InitApp` method in the provided Go code snippet is responsible for initializing the application
// by performing the following tasks:
func (app *App) InitApp() *models.Session {
	fmt.Println("ok App")
	if app.Session != nil {
		app.InitProfile()
		app.SetLogout()
		utils.GetData("getposts", func(r *http.Response, err error) {
			if err != nil {
				fmt.Println("Erreur chargement")
			}
			var postsUI []models.PostUI
			utils.DecodeResponse(r, &postsUI)
			for _, p := range postsUI {
				directive.AllPost = append(directive.AllPost, p)
				go func(p models.PostUI) {
					directive.InitLike(&p.Post)
					directive.InitDislike(&p.Post)
					directive.InitComments(&p.Post)
				}(p)
			}
		})
	} else {
		err := app.AppCommentInit()
		if err != nil {
			fmt.Println(err)
		}
	}
	return app.Session
}

// The `InitProfile` method in the provided Go code snippet is responsible for initializing the user
// profile within the application. Here's a breakdown of what this method does:
func (app *App) InitProfile() {
	user := app.Session.User
	render := utils.ParseComponent(components.Profil_component, user)
	div := directive.GetElement("profil")
	div.Set("innerHTML", render)
}

func (app *App) AppCommentInit() error {
	errChan := make(chan error)
	utils.GetData("getposts", func(r *http.Response, err error) {
		//defer fmt.Println(err)
		var postsUI []models.PostUI
		utils.DecodeResponse(r, &postsUI)
		for _, p := range postsUI {
			if err = directive.InitComments(&p.Post); err != nil {
				errChan <- err
			}
		}
	})
	return <-errChan
}

// The `SetLogout` method in the provided Go code snippet is responsible for setting up a click event
// listener on a button element with the ID "btn-deconnecte". Here's a breakdown of what this method
// does:
func (app *App) SetLogout() {
	btn := directive.GetElement("btn-deconnecte")
	btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		utils.SendData(app.Session, "logout", func(r *http.Response, err error) {
			defer fmt.Println(err)
			js.Global().Get("window").Get("location").Set("href", "/")
		})
		return nil
	}))
}
