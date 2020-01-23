package usecase

import (
	forum "github.com/Andronovdima/tpark-db-forum/internal/app/forum/repository"
	user "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type ForumUsecase struct {
	forumRep *forum.ForumRepository
	userUC  *user.UserUsecase
}

func NewForumUsecase(fr *forum.ForumRepository, ur *user.UserUsecase) *ForumUsecase {
	ForumUsecase := &ForumUsecase{
		forumRep: fr,
		userUC:  ur,
	}
	return ForumUsecase
}

func (u *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	rerr := new(models.HttpError)

	isExistUser := u.userUC.IsExistUser(forum.User)
	if !isExistUser {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "cant find user with this nickname"
		return nil, rerr
	}
	forum.User = u.userUC.GetRightNickname(forum.User)

	isExist := u.forumRep.IsExist(forum.Slug)
	if isExist {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "this forum is exist with that slug"

		fr , err := u.forumRep.Find(forum.Slug)
		if err != nil {
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, rerr
		}

		return fr , rerr
	}

	err := u.forumRep.Create(forum)
	if err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	return forum, nil
}


func (u *ForumUsecase) IsExist(slug string) bool {
	return u.forumRep.IsExist(slug)
}

func (u *ForumUsecase) Find(slug string) (*models.Forum, error) {
	rerr := new(models.HttpError)

	isExist := u.forumRep.IsExist(slug)
	if !isExist {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = " forum doesnt exist with that slug"
		return nil , rerr
	}

	fr , err := u.forumRep.Find(slug)
	if err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil , rerr
	}

	return fr , nil
}


func (u *ForumUsecase) AddThread(slug string) error {
	return u.forumRep.AddThread(slug)
}


func (u *ForumUsecase) AddPost(slug string, count int64) error {
	return u.forumRep.AddPost(slug, count)
}

func (u *ForumUsecase) Status() int64 {
	return u.forumRep.Status()
}

func (u *ForumUsecase) GetRightSlug(slug string) string {
	return u.forumRep.GetRightSlug(slug)
}
