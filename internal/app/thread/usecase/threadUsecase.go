package usecase

import (
	forum "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	thread "github.com/Andronovdima/tpark-db-forum/internal/app/thread/repository"
	user "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type ThreadUsecase struct {
	ThreadRep *thread.ThreadRepository
	UserUC    *user.UserUsecase
	ForumUC   *forum.ForumUsecase
}

func NewUserUsecase(tr *thread.ThreadRepository, uc *user.UserUsecase, fc *forum.ForumUsecase) *ThreadUsecase {
	threadUsecase := &ThreadUsecase{
		ThreadRep: tr,
		UserUC:    uc,
		ForumUC:   fc,
	}
	return threadUsecase
}

func (tu *ThreadUsecase) CreateThread(th *models.Thread, slug string) (*models.Thread, error) {
	rerr := new(models.HttpError)

	isExistForum := tu.ForumUC.IsExist(slug)
	if !isExistForum {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Cant find forum with slug #" + slug
		return nil, rerr
	}

	isExistUser := tu.UserUC.IsExistUser(th.Author)
	if !isExistUser {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Cant find user with author # " + th.Author
		return nil, rerr
	}

	isExistThread := tu.ThreadRep.IsExist(th.Slug)
	if isExistThread {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "That slug already exists"
		th , err := tu.ThreadRep.Find(th.Slug)
		if err != nil {
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, err
		}
		return th, rerr
	}

	err := tu.ThreadRep.Create(th, slug)
	if err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
	}

	return th, nil
}
