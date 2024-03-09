package storage

import (
	"forum/Api/models"
	"forum/orm"
	"forum/utils"
)

type Storage struct {
	Gorm *orm.ORM
}

func NewStorage() *Storage {
	return &Storage{
		Gorm: utils.OrmInit(),
	}
}

func (s *Storage) Insert(model interface{}) error {
	return s.Gorm.Insert(model)
}

func (s *Storage) SelectAll(model interface{}) interface{} {
	_, table := orm.InitTable(model)
	fields := table.GetFieldName()
	return s.Gorm.Scan(model, fields...)
}

func (s *Storage) Select(model interface{}, byField string, value interface{}) interface{} {
	s.Gorm.Custom.Where(byField, value)
	_, table := orm.InitTable(model)
	fields := table.GetFieldName()
	result := s.Gorm.Scan(model, fields...)
	s.Gorm.Custom.Clear()
	return result
}

func (s *Storage) UpdateReaction(model interface{}, id int, column string, newValue interface{}) {
	s.Gorm.SetModel("Id", id, model).UpdateField(newValue, column).Update(s.Gorm.Db)
}

func (s *Storage) UpdateUserReaction(reaction string, post_id, user_id int, table string) error {
	query := "UPDATE "+table+" SET Value = ? WHERE PostId = ? AND User_id = ?"
	_, err := s.Gorm.Db.Exec(query, reaction, post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteExipireSession(user_id int) {
	s.Gorm.Delete(models.SessionDb{}, "User_Id", user_id)
}


