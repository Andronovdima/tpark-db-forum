package usecase

import (
	"github.com/Andronovdima/tpark-db-forum/internal/app/forum"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
)

type ForumUsecase struct {
	forumRep *forum.Repository

}

func NewForumUsecase(fr *forum.Repository) *ForumUsecase {
	ForumUsecase := &ForumUsecase{
		forumRep: fr,
	}
	return ForumUsecase
}

func (u *ForumUsecase) Create(forum *models.Forum) (*models.Forum, error) {
	isExist := u.forumRep.IsExist(forum.Slug)
	if isExist {
		code := 409
		return u.forumRep.GetForum(forum.Slug), nil
	}

	err := forum.Check()
	if err != nil {
		return nil, err
	}

	isExistUser := u.userRep.IsExist(forum.User)
	if !isExistUser {
		code := 404
		return nil, err
	}

	fr ,err := u.forumRep.Create(forum)
	return fr, nil
}
