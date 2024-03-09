package utils

import (
	"encoding/json"
	"forum/Api/models"
	"forum/orm"
	"net/http"
)

// The function `RespondWithJSON` writes JSON data to an HTTP response with the specified status code.
func RespondWithJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// The `OrmInit` function initializes a new ORM instance with a specified database name in Go.
func OrmInit() *orm.ORM {
	gorm := orm.NewORM()
	gorm.InitDB("forumDB.db")
	return gorm
}

// The CreateDatabase function initializes a database connection using GORM and performs auto migration
// for specified models.
func CreateDatabase() {
	gorm := orm.NewORM()
	gorm.InitDB("forumDB.db")
	gorm.AutoMigrate(
		models.User{},
		models.Post{},
		models.ReactionPost{},
		models.ReactionComment{},
		models.Comments{},
		models.Categories{},
		models.SessionDb{},
	)
}
