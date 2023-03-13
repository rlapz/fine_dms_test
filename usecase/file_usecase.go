package usecase

import (
	"errors"
	"fmt"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/utils"
)

var (
	ErrInvalidFileData = errors.New("invalid file data")
	ErrInvalidUserID   = errors.New("invalid user id")
)

type file struct {
	fileRepo repo.FileRepo
}

func NewFileUsecase(fileRepo repo.FileRepo) FileUsecase {
	return &file{fileRepo: fileRepo}
}

func (self *file) GetFilesByUserId(id int) ([]model.File, error) {
	files, err := self.fileRepo.SelectAllByUserId(id)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return nil, ErrUsecaseNoData
		}
		return nil, err
	}

	return files, nil
}

func (self *file) UpdateFile(id int, path, ext string) error {
	if id == 0 || path == "" || ext == "" {
		return ErrInvalidFileData
	}

	files, err := self.fileRepo.SelectAllByUserId(id)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return ErrUsecaseNoData
		}
		return err
	}

	if len(files) == 0 {
		return ErrUsecaseNoData
	}

	file := &files[0]
	file.Path = path
	file.Ext = ext
	file.UpdatedAt = time.Now()

	if err := self.fileRepo.Update(file); err != nil {
		return err
	}

	return nil
}

func (self *file) DeleteFile(userId int, fileId int, token string, secret []byte) error {
	userIdFromToken, err := utils.ValidateToken(token, secret)
	if err != nil {
		return ErrUsecaseInvalidAuth
	}

	if userIdFromToken != fmt.Sprintf("%d", userId) {
		return ErrUsecaseInvalidAuth
	}

	files, err := self.fileRepo.SelectAllByUserId(fileId)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return ErrUsecaseNoData
		}
		return err
	}

	var file *model.File
	for _, f := range files {
		if f.ID == fileId {
			file = &f
			break
		}
	}

	if file == nil {
		return ErrUsecaseNoData
	}

	if file.User.ID != userId {
		return ErrUsecaseInvalidAuth
	}

	// soft delete
	if err := self.fileRepo.Delete(file.ID); err != nil {
		return err
	}

	return nil
}

func (self *file) SearchFilesByUserId(id int, query string) ([]model.File, error) {
	if id == 0 {
		return nil, ErrInvalidUserID
	}

	files, err := self.fileRepo.SelectAllByUserId(id)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return nil, ErrUsecaseNoData
		}
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrUsecaseNoData
	}

	// filter files based on search query
	var filteredFiles []model.File
	for _, f := range files {
		if utils.StringContainsIgnoreCase(f.Path, query) || utils.StringContainsIgnoreCase(f.Ext, query) {
			filteredFiles = append(filteredFiles, f)
		}
	}

	return filteredFiles, nil
}
