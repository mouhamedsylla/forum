package directive

import (
	"fmt"
	"forum/Api/models"
	"forum/App/internal/utils"
	"net/http"
	"strconv"
	"syscall/js"
)

type FilterDirective struct {
	element          js.Value
	SelectCategories map[string]bool
}

func (f *FilterDirective) Selector() string {
	return "[filterForm]"
}
func (f *FilterDirective) SetElement(el js.Value) {
	f.element = el
}

func (f *FilterDirective) Init() {
	f.element.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.FilterBy()
		return nil
	}))
}

func (f *FilterDirective) FilterAll() {
	utils.CheckFinishSession(Session.TimeOut)
	selected := js.Global().Get("document").Call("querySelector", "input[name='created']:checked").Get("value").String()
	if selected == "category" {
		f.InitSelectedCategoryFilter()
		FilterByCategorie(f.SelectCategories)
	}
	if selected == "created-post" {
		FilterByCreated()
	}

	if selected == "post-liked" {
		FilterByLikedPost()
	}

	if selected == "none" {
		for _, postUI := range AllPost {
			div := GetElement("post-" + strconv.Itoa(postUI.Post.Id))
			div.Get("style").Set("display", "block")
		}
	}
}

func (f *FilterDirective) InitSelectedCategoryFilter() {
	f.SelectCategories = map[string]bool{}
	categoryfilter := js.Global().Get("document").Call("querySelectorAll", ".filtercategory")
	for i := 0; i < categoryfilter.Length(); i++ {
		element := categoryfilter.Index(i)
		if element.Get("checked").Bool() {
			f.SelectCategories[element.Get("value").String()] = true
		}
	}
}

func (f *FilterDirective) FilterBy() {
	utils.GetData("getposts", func(r *http.Response, err error) {
		var postUI []models.PostUI
		utils.DecodeResponse(r, &postUI)
		fmt.Println("AllPost: ", AllPost)
		AllPost = postUI
		if Session != nil {
			f.FilterAll()
		} else {
			f.InitSelectedCategoryFilter()
			var tab []string
			for v := range f.SelectCategories {
				tab = append(tab, v)
			}
			if utils.IsValideCategories(tab) {
				FilterByCategorie(f.SelectCategories)
			}
		}
	})
}

func FilterByCategorie(SelectCategories map[string]bool) {
	for _, postUI := range AllPost {
		Exist := false
		for _, v := range postUI.Categories {
			_, ok := SelectCategories[v.Name]
			if ok {
				Exist = true
				break
			}
		}

		div := GetElement("post-" + strconv.Itoa(postUI.Post.Id))
		if !Exist && len(SelectCategories) != 0 {
			div.Get("style").Set("display", "none")
		} else {
			div.Get("style").Set("display", "block")
		}
	}
}

func FilterByLikedPost() {
	utils.GetData("getlikedpost", func(r *http.Response, err error) {
		var postLiked []models.PostUI
		utils.DecodeResponse(r, &postLiked)
		fmt.Println("post liké: ", postLiked)
		for _, pAll := range AllPost {
			div := GetElement("post-" + strconv.Itoa(pAll.Post.Id))
			div.Get("style").Set("display", "none")
		}

		for _, p := range postLiked {
			div := GetElement("post-" + strconv.Itoa(p.Post.Id))
			div.Get("style").Set("display", "block")
		}


		// for _, p := range AllPost {
		// 	go func(p models.PostUI) {
		// 		fmt.Println("id: ", p.Post.Id)
		// 	div := GetElement("post-" + strconv.Itoa(p.Post.Id))
		// 	if len(postLiked) == 0 {
		// 		fmt.Println("ok ok")
		// 		div.Get("style").Set("display", "none")
		// 	}
		// 	for _, v := range postLiked {
		// 		if p.Post.Id != v.Post.Id {
		// 			div.Get("style").Set("display", "none")
		// 		} else {
		// 			div.Get("style").Set("display", "block")	
		// 		}
		// 	}
		// 	}(p)
		// }
	})
}

func FilterByCreated() {
	user := Session.User
	for _, postUI := range AllPost {
		div := GetElement("post-" + strconv.Itoa(postUI.Post.Id))
		if postUI.Post.User_id == user.Id {
			div.Get("style").Set("display", "block")
		} else {
			div.Get("style").Set("display", "none")
		}
	}
}
