package manager

import (
	"enigmacamp.com/fine_dms/usecase"
)

type UsecaseManager interface {
	UserUsecase() usecase.UserUsecase
	TagsUsecase() usecase.TagsUsecase
	FileUseCase() usecase.FileUsecase
	// Add other usecase below
}

type usecaseManager struct {
	repoMgr RepoManager
}

func NewUsecaseManager(repoMgr RepoManager) UsecaseManager {
	return &usecaseManager{
		repoMgr: repoMgr,
	}
}

func (self *usecaseManager) UserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(self.repoMgr.UserRepo())
}

func (self *usecaseManager) TagsUsecase() usecase.TagsUsecase {
	return usecase.NewTagsUsecase(self.repoMgr.TagsRepo())
}

func (self *usecaseManager) FileUseCase() usecase.FileUsecase {
	return usecase.NewFileUsecase(self.repoMgr.FileRepo())
}
