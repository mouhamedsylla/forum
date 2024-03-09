package directive

import (
	"syscall/js"
)

// The above type defines an interface for directives in Go.
// @property Init - The Init method is used to initialize the directive. It can be used to set up any
// necessary state or perform any initializations that are required for the directive to function
// properly.
// @property {string} Selector - The `Selector` method returns a string that represents the CSS
// selector used to select the element(s) that the directive will be applied to.
// @property SetElement - SetElement is a method that sets the element on which the directive is
// applied. It takes a js.Value parameter, which represents the JavaScript object that corresponds to
// the element.
type Directive interface {
	Init()
	Selector() string
	SetElement(js.Value)
}
