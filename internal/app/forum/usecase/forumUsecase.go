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

func (u *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, error) {
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


func (u *ForumUsecase) CreateThread(th *models.Thread, slug string) (*models.Thread, error) {
	err := new(models.HttpError)

	isExistForum := u.forumRep.IsExist(slug)
	if !isExistForum {
		err.StatusCode = 404
		err.StringErr = "Cant find forum with slug #" + slug
		return nil, err
	}

	isExistUser := u.userRep.IsExist(th.Author)
	if !isExistUser {
		err.StatusCode = 404
		err.StringErr = "Cant find user with author # " + th.Author
		return nil, err
	}

	isExistSlug := u.threadRep.IsExist(th)
	if isExistSlug {
		err := new(models.HttpError)
		err.StatusCode = 409
		err.StringErr = "That slug already exists"
		return u.threadRep.Give(th), err
	}

	cerr := u.threadRep.Create(th)
	if cerr != nil {
		err.StatusCode = 500
		err.StringErr = cerr.Error()
	}

	return th, nil
}
