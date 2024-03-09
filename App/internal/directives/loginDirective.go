package directive

import (
	"errors"
	"forum/Api/models"
	"forum/App/internal/utils"
	"net/http"
	"syscall/js"
	"time"
)

var Session *models.Session

type LoginDirective struct {
	element js.Value
	user    models.User
}

func (l *LoginDirective) Selector() string {
	return "[login]"
}

func (l *LoginDirective) SetElement(el js.Value) {
	l.element = el
}

func (l *LoginDirective) Init() {
	l.element.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")
		go func() {
			err := l.GetUserLogin()
			if err != nil {
				ErrorConnect(err.Error())
				return
			}
			l.Login()
		}()
		return nil
	}))
}

func (l *LoginDirective) Login() {
	utils.SendData(l.user, "login", func(r *http.Response, err error) {
		var message *models.Message
		utils.DecodeResponse(r, &message)
		if len(message.Error) != 0 {
			ErrorConnect(message.Error)
		} else {
			js.Global().Get("window").Get("location").Set("href", "/home")
		}
	})

}

func (l *LoginDirective) GetUserLogin() error {
	l.user.Email = GetElement("userlog").Get("value").String()
	l.user.Password = GetElement("userpass").Get("value").String()
	if !utils.NotNullEntry(l.user.Email, l.user.Password) {
		return errors.New("incomplete entries")
	}
	return nil
}

//time after n'est pas bloquant contrairement à time.Spleep 
func ErrorConnect(err string) {
	divErr := js.Global().Get("document").Call("getElementById", "errorLogin")
	divErr.Set("innerHTML", "<p>"+err+"</p>")
	divErr.Get("style").Set("display", "block")
	<-time.After(1000 * time.Millisecond)
	divErr.Get("style").Set("display", "none")
}
