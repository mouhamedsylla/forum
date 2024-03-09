package directive

import (
	"errors"
	"forum/Api/models"
	"forum/App/internal/utils"
	"net/http"
	"syscall/js"
)

type RegisterDirective struct {
	element js.Value
	user    models.User
	inputErr InputError
}

func (l *RegisterDirective) Selector() string {
	return "[register]"
}

func (l *RegisterDirective) SetElement(el js.Value) {
	l.element = el
}

func (l *RegisterDirective) Init() {
	l.element.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")
		l.inputErr = InputError{}
		l.inputErr = ErrorInput
		go func() {
			err := l.GetUserRegistration()
			if err != nil {
				ErrorConnect(err.Error())
				return
			}
			if l.inputErr.message != nil {
				ErrorConnect(l.inputErr.message.Error())
				return
			}
			l.Registration()
		}()
		return nil
	}))
}

func (l *RegisterDirective) Registration() {
	utils.SendData(l.user, "register", func(r *http.Response, err error) {
		var message models.Message
		utils.DecodeResponse(r, &message)
		if len(message.Error) != 0 {
			ErrorConnect(message.Error)
		} else {
			js.Global().Get("window").Get("location").Set("href", "/connect")
		}
	})
}

func (l *RegisterDirective) GetUserRegistration() error {
	l.user.Username = GetElement("user").Get("value").String()
	l.user.Email = GetElement("email").Get("value").String()
	l.user.Password = GetElement("password").Get("value").String()
	if !utils.NotNullEntry(l.user.Username, l.user.Email, l.user.Password) {
		return errors.New("incomplete entries")
	}
	return nil
}

func GetElement(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}
