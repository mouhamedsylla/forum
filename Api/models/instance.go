package models

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"

	"golang.org/x/crypto/bcrypt"
)

func CryptPassword(user *User) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hashPassword)
}

func New(r *http.Request, model interface{}) (interface{}, int, error) {
	newStruct := reflect.New(reflect.TypeOf(model)).Interface()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err = json.Unmarshal(data, newStruct); err != nil {
		return nil, http.StatusBadRequest, err
	}
	return newStruct, http.StatusOK, nil
}
