package main

import (
	"fmt"
	"forum/App/internal/boot"
	directive "forum/App/internal/directives"
)

func main() {
	fmt.Println("WebAssembly from Golang")

	app := boot.NewApp()
	app.Directives = []directive.Directive{
		&directive.LoginDirective{},
		&directive.RegisterDirective{},
		&directive.PostDirective{},
		&directive.CommentDirective{},
		&directive.InputValidate{},
		&directive.FilterDirective{},
	}

	app.BootstrapApplication()
	app.InitSession()
	directive.Session = app.InitApp()
	select {}
}
