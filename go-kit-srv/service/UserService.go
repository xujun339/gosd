package service

import (
	"go-kit-srv/service/discovery"
	"strconv"
)

type IUserService interface {
	GetName(userId int) string
}

type UserService struct {

}

func (user UserService) GetName(userId int) string {
	if userId == 101 {
		return "gavin" + strconv.Itoa(discovery.GetServicePort())
	} else {
		return "guest"
	}
}
