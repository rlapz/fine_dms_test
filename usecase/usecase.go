package usecase

import (
	"errors"

	"enigmacamp.com/fine_dms/model"
)

var (
	ErrUsecaseInternal       = errors.New("internal server error")
	ErrUsecaseNoData         = errors.New("no data")
	ErrUsecaseEmptyEmail     = errors.New("`email` cannot be empty")
	ErrUsecaseEmptyUsername  = errors.New("`username` cannot be empty")
	ErrUsecaseEmptyPassword  = errors.New("`password` cannot be empty")
	ErrUsecaseEmptyFname     = errors.New("`first_name` cannot be empty")
	ErrUsecaseExistsUsername = errors.New("`username` already exists")
	ErrUsecaseExistsEmail    = errors.New("`email` already exists")
	ErrUsecaseInvalidEmail   = errors.New("`email` invalid format")
	ErrUsecaseInvalidAuth    = errors.New("`username` or `password` wrong")
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

type FileUsecase interface {
	GetFilesByUserId(id int) ([]model.File, error)
	UpdateFile(id int, path, ext string) error
	DeleteFile(userId int, fileId int, token string, secret []byte) error
	SearchFilesByUserId(id int, query string) ([]model.File, error)
}
