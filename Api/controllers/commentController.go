package controllers

import (
	"fmt"
	"forum/Api/models"
	"forum/utils"
	"net/http"
	"strconv"
)

var CommentUI []models.Comments
var Comment_reaction = make(map[string]string)

func (c *Controllers) Comment() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.Comments{})
		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		comment := result.(*models.Comments)
		err = c.Storage.Insert(*comment)
		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, http.StatusInternalServerError)
			return
		}
		c.Storage.Gorm.Custom.OrderBy("Id", 1).Limit(1)
		dbComment := c.Storage.Gorm.Scan(models.Comments{}, "Id").([]models.Comments)[0]
		c.Storage.Gorm.Custom.Clear()
		utils.RespondWithJSON(w, models.CommentId{Id: dbComment.Id}, http.StatusOK)
	})
}

func (c *Controllers) GetComments() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.Post{})
		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		post := result.(*models.Post)
		allComments := c.Storage.Select(models.Comments{}, "Post_Id", post.Id).([]models.Comments)
		utils.RespondWithJSON(w, allComments, http.StatusOK)
	})
}

func (c *Controllers) CommentUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, models.Comments{})
		if err != nil {
			message.Error = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		comment := result.(*models.Comments)
		rslt, ok := c.Storage.Select(models.User{}, "Id", comment.User_id).([]models.User)
		if !ok {
			message.Error = "user not found"
			utils.RespondWithJSON(w, message, http.StatusNotFound)
			return
		}
		user := rslt[0]
		utils.RespondWithJSON(w, user, http.StatusOK)
	})
}

func (c *Controllers) CommentReaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("user").(*models.Session)
		var message models.Message
		result, status, err := models.New(r, models.Action{})
		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		action := result.(*models.Action)
		comment := c.Storage.Select(models.Comments{}, "Id", action.Id).([]models.Comments)[0]
		key := strconv.Itoa(action.Id) + "-" + strconv.Itoa(session.User.Id)
		rct, ok := Comment_reaction[key]
		if !ok {
			reaction := models.ReactionComment{
				Value:   action.Reaction,
				CommentId:  action.Id,
				User_id: session.User.Id,
			}
			c.Storage.Insert(reaction)
			Comment_reaction[key] = action.Reaction
			if action.Reaction == "Like" {
				c.Storage.UpdateReaction(models.Comments{}, action.Id, action.Reaction, comment.Like+1)
			} else {
				c.Storage.UpdateReaction(models.Comments{}, action.Id, action.Reaction, comment.Dislike+1)
			}

		} else {
			if rct == "Like" {
				if action.Reaction == "Dislike" {
					c.Storage.UpdateReaction(models.Comments{}, action.Id, action.Reaction, comment.Dislike+1)
					if comment.Like > 0 {
						comment.Like = comment.Like - 1
					}
					c.Storage.UpdateReaction(models.Comments{}, action.Id, "Like", comment.Like)
				}
			} else {
				if action.Reaction == "Like" {
					c.Storage.UpdateReaction(models.Comments{}, action.Id, action.Reaction, comment.Like+1)
					if comment.Dislike > 0 {
						comment.Dislike = comment.Dislike - 1
					}
					c.Storage.UpdateReaction(models.Comments{}, action.Id, "Dislike", comment.Dislike)
				}
			}
		}
		Comment_reaction[key] = action.Reaction
		err = c.Storage.UpdateUserReaction(action.Reaction, action.Id, session.User.Id, "ReactionComment")
		if err != nil {
			fmt.Println(err)
		}
		comment = c.Storage.Select(models.Comments{}, "Id", action.Id).([]models.Comments)[0]
		utils.RespondWithJSON(w, comment, http.StatusOK)
	})
}
