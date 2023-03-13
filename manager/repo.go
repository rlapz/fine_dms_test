package manager

import (
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/repo/psql"
)

type RepoManager interface {
	UserRepo() repo.UserRepo
	TagsRepo() repo.TagsRepo
	FileRepo() repo.FileRepo
	// Add other repo below
}

type repoManager struct {
	infraMgr InfraManager
}

func NewRepoManager(infr InfraManager) RepoManager {
	return &repoManager{
		infraMgr: infr,
	}
}

func (self *repoManager) UserRepo() repo.UserRepo {
	return psql.NewPsqlUserRepo(self.infraMgr.GetDB())
}

func (self *repoManager) TagsRepo() repo.TagsRepo {
	return psql.NewPsqlTagsRepo(self.infraMgr.GetDB())
}

func (self *repoManager) FileRepo() repo.FileRepo {
	return psql.NewPsqlFileRepo(self.infraMgr.GetDB())
}
