package usecase

import (
	"errors"

	"enigmacamp.com/fine_dms/model"
)

var (
	ErrUsecaseInternal = errors.New("internal server error")
	ErrUsecaseNoData   = errors.New("no data")
)

type TagsUsecase interface {
	GetAll() ([]model.Tags, error)
}

type UserUsecase interface {
	GetAll() ([]model.User, error)
	GetById(id int) (*model.User, error)
	GetByUsername(uname string) (*model.User, error)
	Add(user *model.User) error
	Edit(user *model.User) error
	Del(id int) error
	AuthenticateUser(string, string) (int64, error)
}
