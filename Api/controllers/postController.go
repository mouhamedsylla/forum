package controllers

import (
	"fmt"
	"forum/Api/models"
	"forum/utils"
	"net/http"
	"strconv"
)

var Post_reaction = make(map[string]string)
var PostsUI []*models.PostUI

func (c *Controllers) CreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		// New prend en paramètre une requête et un modèle de données pour extraire
		// les données contenu dans le requête et retourne un interface
		result, status, err := models.New(r, models.Post{})
		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		// l'interface ainsi retourné sera convertie dans le modèle de données spécifié 
		// dés le depart
		post := result.(*models.Post)
		post.Title = "helloworld"

		// stockage du post récupèré dans la base de données
		err = c.Storage.Insert(*post)
		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, http.StatusBadRequest)
			return
		}

		// ce requête récupère le dernier post 
		c.Storage.Gorm.Custom.OrderBy("Id", 1).Limit(1)
		dbPost := c.Storage.Gorm.Scan(models.Post{}, "Id").([]models.Post)[0]
		c.Storage.Gorm.Custom.Clear()
		// Renvoie le post ainsi stocker avec son l'Id à l'utilisateur 
		utils.RespondWithJSON(w, models.PostID{Id: dbPost.Id}, http.StatusOK)
	})
}

func (c *Controllers) GetPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, PostsUI, http.StatusOK)
	})
}

func (c *Controllers) GetLikedPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("user").(*models.Session)
		user := session.User
		c.Storage.Gorm.Custom.Where("Value", "Like").And("User_id", user.Id)
		rslt := c.Storage.Gorm.Scan(models.ReactionPost{}, "PostId").([]models.ReactionPost)
		c.Storage.Gorm.Custom.Clear()
		if len(rslt) == 0{
			fmt.Println("no liked post")

			return
		}
		var postsUI []models.PostUI
		fmt.Println("rslt: ", rslt)
		for _, rp := range rslt {
			for _, pst := range PostsUI {
				if rp.PostId == pst.Post.Id {
					postsUI = append(postsUI, *pst)
				}
			}
		}
		fmt.Println("postUI: ", postsUI)
		utils.RespondWithJSON(w, postsUI, http.StatusOK)
	})
}

func (c *Controllers) PostInfos() {
	Posts := c.Storage.SelectAll(models.Post{}).([]models.Post)
	for _, p := range Posts {
		var categories []models.Categories
		attemp := c.Storage.Select(models.Categories{}, "Id_Post", p.Id)
		if attemp != nil {
			categories = attemp.([]models.Categories)
		}
		user := c.Storage.Select(models.User{}, "Id", p.User_id).([]models.User)[0]
		PostsUI = append(PostsUI, &models.PostUI{
			User:       user,
			Post:       p,
			Categories: categories,
		})
	}
}

func (c *Controllers) AddCategories() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message models.Message
		result, status, err := models.New(r, []models.Categories{})
		if err != nil {
			message.Message = err.Error()
			utils.RespondWithJSON(w, message, status)
			return
		}
		categories := result.(*[]models.Categories)
		if len(*categories) == 0 {
			message.Error = "you must choice one category"
			utils.RespondWithJSON(w, message, http.StatusBadRequest)
			return
		}
		for _, category := range *categories {
			err = c.Storage.Insert(category)
			if err != nil {
				message.Message = err.Error()
				utils.RespondWithJSON(w, message, http.StatusBadRequest)
				return
			}
		}
		AddPostUI(*categories, c)
		message.Message = "category added successful"
		utils.RespondWithJSON(w, message, http.StatusOK)
	})
}

func (c *Controllers) Reaction() http.Handler {
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
		post := c.Storage.Select(models.Post{}, "Id", action.Id).([]models.Post)[0]
		key := strconv.Itoa(action.Id) + "-" + strconv.Itoa(session.User.Id)
		rct, ok := Post_reaction[key]
		if !ok {
			reaction := models.ReactionPost{
				Value:   action.Reaction,
				PostId:  action.Id,
				User_id: session.User.Id,
			}
			c.Storage.Insert(reaction)
			Post_reaction[key] = action.Reaction
			if action.Reaction == "Like" {
				c.Storage.UpdateReaction(models.Post{}, action.Id, action.Reaction, post.Like+1)
			} else {
				c.Storage.UpdateReaction(models.Post{}, action.Id, action.Reaction, post.Dislike+1)
			}

		} else {
			if rct == "Like" {
				if action.Reaction == "Dislike" {
					c.Storage.UpdateReaction(models.Post{}, action.Id, action.Reaction, post.Dislike+1)
					if post.Like > 0 {
						post.Like = post.Like - 1
					}
					c.Storage.UpdateReaction(models.Post{}, action.Id, "Like", post.Like)
				}
			} else {
				if action.Reaction == "Like" {
					c.Storage.UpdateReaction(models.Post{}, action.Id, action.Reaction, post.Like+1)
					if post.Dislike > 0 {
						post.Dislike = post.Dislike - 1
					}
					c.Storage.UpdateReaction(models.Post{}, action.Id, "Dislike", post.Dislike)
				}
			}
		}
		Post_reaction[key] = action.Reaction
		err = c.Storage.UpdateUserReaction(action.Reaction, action.Id, session.User.Id, "ReactionPost")
		if err != nil {
			fmt.Println(err)
		}
		post = c.Storage.Select(models.Post{}, "Id", action.Id).([]models.Post)[0]
		UpdateReaction(post)
		utils.RespondWithJSON(w, post, http.StatusOK)
	})
}

func AddPostUI(categories []models.Categories, c *Controllers) {
	idPost := categories[0].Id_Post
	post := c.Storage.Select(models.Post{}, "Id", idPost).([]models.Post)[0]
	user := c.Storage.Select(models.User{}, "Id", post.User_id).([]models.User)[0]
	newPostUI := &models.PostUI{
		User:       user,
		Post:       post,
		Categories: categories,
	}
	PostsUI = append(PostsUI, newPostUI)
}

func UpdateReaction(post models.Post) {
	for _, pUI := range PostsUI {
		if pUI.Post.Id == post.Id {
			pUI.Post.Like = post.Like
			pUI.Post.Dislike = post.Dislike
		}
	}
}
