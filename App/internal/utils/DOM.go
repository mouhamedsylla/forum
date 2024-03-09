package utils

import "syscall/js"

type HtmlElement struct {
	Element    js.Value
	TagName    js.Value
	Id         js.Value
	ClassList  []js.Value
	Attributes map[string]string
	Content    string
	Style      map[string]string
}

func (htE *HtmlElement) GetTagName() js.Value {
	return htE.Element.Get("tagName")
}

func (htE *HtmlElement) GetId() js.Value {
	return htE.Element.Get("id")
}

func (htE *HtmlElement) GetClassList() []js.Value {
	classList := htE.Element.Get("classlist")
	classList.Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		htE.ClassList = append(htE.ClassList, args[0])
		return nil
	}))
	return htE.ClassList
}

func (htE *HtmlElement) GetAttributes() {

}
