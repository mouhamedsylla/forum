package directive

import (
	"errors"
	"regexp"
	"syscall/js"
)

type ValidateFunc func() (bool, error)

type InputError struct {
	inputName string
	message error
}

var (
	ErrorInput = InputError{}
	Valide int
)


type InputValidate struct {
	element js.Value
}

func (input *InputValidate) Selector() string {
	return "[inputValidate]"
}

func (input *InputValidate) SetElement(el js.Value) {
	input.element = el
}

func (input *InputValidate) Init() {
	var useFunc ValidateFunc
	if input.HaveAttribut("username") {
		useFunc = input.ValideUsername
	}
	if input.HaveAttribut("email") {
		useFunc = input.ValideEmail
	}
	if input.HaveAttribut("password") {
		useFunc = input.ValidePassword
	}
	input.element.Call("addEventListener", "input", input.InputCallback(useFunc))
}

func (input *InputValidate) InputCallback(f ValidateFunc) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		result, err := f()
		if err != nil || !result {
			input.element.Get("style").Set("border-bottom", "2px solid red")
			Valide++
		}

		if err != nil {
			ErrorInput.inputName = this.Get("id").String()
			ErrorInput.message = errors.New(err.Error())
			
		} else if !result {
			ErrorInput.inputName = this.Get("id").String()
			ErrorInput.message = errors.New( "invalide " + ErrorInput.inputName )
		}
		if result {
			input.element.Get("style").Set("border-bottom", "2px solid green")
			Valide--
			ErrorInput = InputError{}
		}
		return nil
	})
}

func (input *InputValidate) ValideUsername() (bool, error) {
	input.element = GetElement("user")
	username := input.element.Get("value").String()
	if len(username) > 20 {
		input.element.Set("value", username[:19])
		return false, errors.New("size limit exceeded")
	}
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return usernameRegex.MatchString(username), nil
}

func (input *InputValidate) ValidePassword() (bool, error) {
	input.element = GetElement("password")
	password := input.element.Get("value").String()
	if len(password) < 8 || len(password) > 15 {
		if len(password) > 15 {
			input.element.Set("value", password[:14])
		}
		return false, errors.New("your password must be between 8 and 15 characters")
	}
	var (
		hasMin     = regexp.MustCompile(`[a-z]`).MatchString
		hasMaj     = regexp.MustCompile(`[A-Z]`).MatchString
		hasDigit   = regexp.MustCompile(`\d`).MatchString
		hasSpecial = regexp.MustCompile(`[\W_]`).MatchString
	)

	return hasMin(password) && hasMaj(password) && hasDigit(password) && hasSpecial(password), nil
}

func (input *InputValidate) ValideEmail() (bool, error) {
	input.element = GetElement("email")
	email := input.element.Get("value").String()
	if len(email) > 40 {
		input.element.Set("value", email[:39])
		return false, errors.New("size limit exceeded")
	}
	var emailRegex = regexp.MustCompile(`^\S+@\S+\.\S+$`)
	return emailRegex.MatchString(email), nil
}

func (input *InputValidate) HaveAttribut(attribut string) bool {
	return input.element.Call("hasAttribute", attribut).Bool()
}
