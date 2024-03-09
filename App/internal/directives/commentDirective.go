package directive

import (
	"fmt"
	"forum/Api/models"
	"forum/App/internal/components"
	"forum/App/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
)

type CommentDirective struct {
	element js.Value
}

func (c *CommentDirective) Selector() string {
	return "[commentaire]"
}

func (c *CommentDirective) SetElement(el js.Value) {
	c.element = el
}

func (c *CommentDirective) Init() {
	c.handleCommentClick(c.element)
}

func (c *CommentDirective) handleCommentClick(el js.Value) {
	el.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		utils.CheckFinishSession(Session.TimeOut)
		document := js.Global().Get("document")
		textarea := document.Call("getElementById", "comment-text")
		texte := textarea.Get("value").String()
		comment := &models.Comments{
			Comment: texte,
			Post_id: PostClick,
			User_id: Session.User.Id,
		}
		if len(strings.TrimSpace(comment.Comment)) != 0 {
			c.BuildComment(comment)
		}
		
		return nil
	}))
}

func (c *CommentDirective) BuildComment(comment *models.Comments) {
	utils.SendData(comment, "comment", func(r *http.Response, err error) {
		defer fmt.Println(err)
		var commentId models.CommentId
		utils.DecodeResponse(r, &commentId)
		comment.Id = commentId.Id
		CreateComment(comment, Session.User)
		InitLikeComment(comment)
		InitDislikeComment(comment)
	})
}

func CreateComment(comment *models.Comments, user *models.User) {
	render := utils.ParseComponent(components.Comment_componnent, comment)
	commentsDiv := GetElement("comments")
	div := js.Global().Get("document").Call("createElement", "div")
	div.Get("classList").Call("add", "feedcmt")
	div.Call("setAttribute", "id", "comment-"+strconv.Itoa(comment.Id))
	div.Set("innerHTML", render)
	commentsDiv.Call("appendChild", div)
	h3 := GetElement("usercomment-" + strconv.Itoa(comment.Id))
	h3.Set("innerText", user.Username)
}

func InitLikeComment(comment *models.Comments) {
	likeElement := GetElement("likedcomment-" + strconv.Itoa(comment.Id))
	likeElement.Call("addEventListener", "click", SetReactionComment("Like", comment.Id))
}

func InitDislikeComment(comment *models.Comments) {
	dislikeElement := GetElement("dislikedcomment-" + strconv.Itoa(comment.Id))
	dislikeElement.Call("addEventListener", "click", SetReactionComment("Dislike", comment.Id))
}


func SetReactionComment(reaction string, commentId int) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := GetElement("reactionComment-" + strconv.Itoa(commentId))
		input.Set("value", reaction)
		action := models.Action{
			Reaction: reaction,
			Id:   commentId,
		}
		fmt.Println("action: ", action)
		utils.SendData(action, "commentreaction", func(r *http.Response, err error) {
			var comment models.Comments
			utils.DecodeResponse(r, &comment)
			if comment.Id != 0 {
				GetElement("likedcomment-"+strconv.Itoa(commentId)).Set("innerText", strconv.Itoa(comment.Like))
				GetElement("dislikedcomment-"+strconv.Itoa(commentId)).Set("innerText", strconv.Itoa(comment.Dislike))
			}
		})
		return nil
	})
}
