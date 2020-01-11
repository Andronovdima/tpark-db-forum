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
	isExist := u.forumRep.IsExist(forum.Slug)
	if isExist {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "this forum is exist with that slug"

		fr , err := u.forumRep.Find(forum.Slug)
		if err != nil {
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, err
		}

		return fr , rerr
	}

	isExistUser := u.userUC.IsExistUser(forum.User)
	if !isExistUser {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "cant find user with this nickname"
		return nil, rerr
	}

	err := u.forumRep.Create(forum)

	return forum, err
}


func (u *ForumUsecase) IsExist(slug string) bool {
	return u.forumRep.IsExist(slug)
}