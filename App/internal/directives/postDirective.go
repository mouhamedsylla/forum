package directive

import (
	"errors"
	"fmt"
	"forum/Api/models"
	"forum/App/internal/components"
	"forum/App/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

var (
	AllPost   []models.PostUI
	PostClick int
)

type PostDirective struct {
	element js.Value
}

func (p *PostDirective) Selector() string {
	return "[post]"
}

func (p *PostDirective) SetElement(el js.Value) {
	p.element = el
}

func (p *PostDirective) Init() {
	fmt.Println("OKK init")
	p.CreatePost(p.element)
}

func (p *PostDirective) CreatePost(el js.Value) {
	el.Call("addEventListener", "click", p.handlePostClick())
}

func (p *PostDirective) handlePostClick() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")
		utils.CheckFinishSession(Session.TimeOut)
		text := js.Global().Get("document").Call("getElementById", "text-post").Get("value").String()

		// Créer le post
		post := &models.Post{
			Content: text,
			User_id: Session.User.Id,
		}
		post.CreatedAt = time.Now().Format("15:04:05")
		// BuildPost fait une requête post pour envoyer les donnees vers le serveur qui le stcoker dans la DB
		err := p.BuildPost(post)
		if err != nil {
			DisplayErrorPost(err)
		}
		return nil
	})
}

func DisplayErrorPost(err error) {
	p := GetElement("create-post-error")
	p.Set("innerText", err.Error())
	p.Set("innerTetx", "")
}

func (p *PostDirective) BuildPost(post *models.Post) error {
	categories := getSelectedCategories()
	if len(categories) == 0 {
		return errors.New("choose at least one category")
	}
	if !utils.IsValideCategories(categories) {
		return errors.New("invalide categories")
	}
	if len(strings.TrimSpace(post.Content)) == 0 {
		return errors.New("text required")
	}
	// Envoyer les données du post au serveur
	utils.SendData(post, "post", func(r *http.Response, err error) {
		if err != nil {
			fmt.Println("Erreur lors de la création du post:", err)
		}
		var postId *models.PostID
		utils.DecodeResponse(r, &postId)
		post.Id = postId.Id

		BuildCategories(postId.Id)
		UpdateUI(categories, post, Session.User)
		GetElement("pupPupModal").Get("style").Set("display", "none")
	})
	return nil
}

func AddPost(post *models.Post, user *models.User) {
	render := utils.ParseComponent(components.Post_component, post)
	postContainer := GetElement("feeds")
	div := js.Global().Get("document").Call("createElement", "div")
	div.Get("classList").Call("add", "feed")
	div.Call("setAttribute", "id", "post-"+strconv.Itoa(post.Id))
	div.Set("innerHTML", render)
	postContainer.Call("appendChild", div)
	h3 := GetElement("name-" + strconv.Itoa(post.Id))
	h3.Set("innerText", user.Username)
}

func BuildCategories(postID int) {
	// Récupérer les catégories sélectionnées
	var post_categories []models.Categories
	categories := getSelectedCategories()
	for _, c := range categories {
		post_categories = append(post_categories, models.Categories{
			Name:    c,
			Id_Post: postID,
		})
	}
	// Envoyer les catégories au serveur
	utils.SendData(post_categories, "category", func(r *http.Response, err error) {
		if err != nil {
			fmt.Println("Erreur lors de l'envoi des catégories:", err)
		}
	})
}

func getSelectedCategories() []string {
	var selectedCategories []string
	document := js.Global().Get("document")
	categoryInputs := document.Call("querySelectorAll", ".category")
	for i := 0; i < categoryInputs.Length(); i++ {
		if categoryInputs.Index(i).Get("checked").Bool() {
			selectedCategories = append(selectedCategories, categoryInputs.Index(i).Get("value").String())
		}
	}
	return selectedCategories
}

func InitLike(post *models.Post) {
	likeElement := GetElement("liked-" + strconv.Itoa(post.Id))
	likeElement.Call("addEventListener", "click", SetReactionPost("Like", post.Id))
}

func InitDislike(post *models.Post) {
	dislikeElement := GetElement("disliked-" + strconv.Itoa(post.Id))
	dislikeElement.Call("addEventListener", "click", SetReactionPost("Dislike", post.Id))
}
func InitComments(post *models.Post) error {
	commentsElement := GetElement("comment-" + strconv.Itoa(post.Id))
	if commentsElement.IsNull() {
		return errors.New("element not found")
	}
	commentsElement.Call("addEventListener", "click", handleCommentButton(post))
	return nil
}

func handleCommentButton(post *models.Post) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		//utils.CheckFinishSession(Session.TimeOut)
		PostClick = post.Id
		fmt.Println("PostClick: ", PostClick)
		//InitUIPost()
		input := GetElement("state-" + strconv.Itoa(post.Id))
		state := input.Get("value").String()
		if state == "close" {
			commentsDiv := GetElement("comments")
			commentsDiv.Set("innerHTML", "")
			BuildCommentPost(post)
			PostDisplay("none")
			GetElement("comments-container").Get("style").Set("display", "block")
			//GetElement("comment-text").Get("style").Set("display", "block")
			input.Set("value", "open")
		}
		if state == "open" {
			PostDisplay("block")
			GetElement("comments-container").Get("style").Set("display", "none")
			//GetElement("comment-text").Get("style").Set("display", "none")
			input.Set("value", "close")
		}
		return nil
	})
}

func PostDisplay(value string) {
	utils.GetData("getposts", func(r *http.Response, err error) {
		var posts []models.PostUI
		utils.DecodeResponse(r, &posts)
		AllPost = posts
		for _, pUI := range AllPost {
			if pUI.Post.Id != PostClick {
				p := GetElement("post-" + strconv.Itoa(pUI.Post.Id))
				p.Get("style").Set("display", value)
			}
		}
	})
}

func BuildCommentPost(post *models.Post) {
	utils.SendData(post, "getcomment", func(r *http.Response, err error) {
		var comments []models.Comments
		utils.DecodeResponse(r, &comments)
		for _, c := range comments {
			CommentUserCreate(c)
		}
	})
}

func CommentUserCreate(comment models.Comments) {
	utils.SendData(comment, "getcommentuser", func(r *http.Response, err error) {
		var user models.User
		utils.DecodeResponse(r, &user)
		CreateComment(&comment, &user)
		InitLikeComment(&comment)
		InitDislikeComment(&comment)
	})
}

func SetReactionPost(reaction string, postId int) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := GetElement("reaction-" + strconv.Itoa(postId))
		input.Set("value", reaction)
		action := models.Action{
			Reaction: reaction,
			Id:       postId,
		}
		utils.SendData(action, "reaction", func(r *http.Response, err error) {
			var post models.Post
			utils.DecodeResponse(r, &post)
			if post.Id != 0 {
				GetElement("liked-"+strconv.Itoa(postId)).Set("innerText", strconv.Itoa(post.Like))
				GetElement("disliked-"+strconv.Itoa(postId)).Set("innerText", strconv.Itoa(post.Dislike))
			}
		})
		return nil
	})
}

func UpdateUI(categories []string, post *models.Post, user *models.User) {
	AddPost(post, user)
	InitLike(post)
	InitDislike(post)
	InitComments(post)
	GetElement("typeCategorie-"+strconv.Itoa(post.Id)).Set("innerText", strings.Join(categories, "    "))
}
