package controllers

import (
	"forum/Api/services"
	"forum/Api/storage"
)

type Controllers struct {
	Auth    *services.AuthService
	Storage *storage.Storage
}

func NewControllers() *Controllers {
	return &Controllers{
		Auth:    services.NewAuthService(),
		Storage: storage.NewStorage(),
	}
}
