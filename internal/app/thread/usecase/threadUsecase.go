package usecase

import (
	forum "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	thread "github.com/Andronovdima/tpark-db-forum/internal/app/thread/repository"
	user "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strconv"
)

type ThreadUsecase struct {
	threadRep *thread.ThreadRepository
	userUC    *user.UserUsecase
	forumUC   *forum.ForumUsecase
}

func NewUserUsecase(tr *thread.ThreadRepository, uc *user.UserUsecase, fc *forum.ForumUsecase) *ThreadUsecase {
	threadUsecase := &ThreadUsecase{
		threadRep: tr,
		userUC:    uc,
		forumUC:   fc,
	}
	return threadUsecase
}

func (tu *ThreadUsecase) CreateThread(th *models.Thread, slug string) (*models.Thread, error) {
	rerr := new(models.HttpError)

	isExistForum := tu.forumUC.IsExist(slug)
	if !isExistForum {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Cant find forum with slug #" + slug
		return nil, rerr
	}
	slug = tu.forumUC.GetRightSlug(slug)

	isExistUser := tu.userUC.IsExistUser(th.Author)
	if !isExistUser {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Cant find user with author # " + th.Author
		return nil, rerr
	}
	th.Author = tu.userUC.GetRightNickname(th.Author)

	isExistThread, id := tu.threadRep.IsExist(th)
	if isExistThread {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "That slug already exists"
		th , err := tu.threadRep.FindByID(id)
		if err != nil {
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, err
		}
		return th, rerr
	}

	err := tu.threadRep.Create(th, slug)
	if err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil,rerr
	}

	_ = tu.forumUC.AddThread(slug)


	return th, nil
}

func (tu *ThreadUsecase) GetThreads(forumSlug string, limit int, created string, desc bool) ([]models.Thread, error) {
	isExist := tu.forumUC.IsExist(forumSlug)
	if !isExist {
		err := new(models.HttpError)
		err.StatusCode = http.StatusNotFound
		err.StringErr = "Forum with this slug doesn't exist"
		return nil, err
	}

	threads , err := tu.threadRep.GetThreads(forumSlug, limit, created, desc)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	return threads, nil
}

//func (tu *ThreadUsecase) IsExist(slug string) bool {
//	 return tu.threadRep.IsExist(slug)
//}

func (tu *ThreadUsecase) IsExistByID(ID int) bool {
	return tu.threadRep.IsExistByID(ID)
}

func (tu *ThreadUsecase) GetIDBySlug(slug string) int {
	return tu.threadRep.GetIDBySlug(slug)
}


func (tu *ThreadUsecase) GetForum(slug string, id int32) (string, error) {
	return tu.threadRep.GetForum(slug, id)
}

func (tu *ThreadUsecase) GetThreadByID(id int32) (*models.Thread, error) {
	currThread , err := tu.threadRep.FindByID(id)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	return currThread, nil
}

func (tu *ThreadUsecase) EditThread(slugID string, thrUpd *models.ThreadUpdate) (*models.Thread, error) {
	ID, err := strconv.Atoi(slugID)
	if err != nil {
		isExist := tu.threadRep.IsExistBySlug(slugID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this id doesnt't exist"
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}
		ID = tu.threadRep.GetIDBySlug(slugID)
	} else {
		isExist := tu.threadRep.IsExistByID(ID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this id doesnt't exist"
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}
	}
	thr, err := tu.threadRep.FindByID(int32(ID))
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	if thrUpd.Message == "" {
		thrUpd.Message = thr.Message
	}

	if thrUpd.Title == "" {
		thrUpd.Title = thr.Title
	}

	err = tu.threadRep.EditThread(int32(ID), thrUpd)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	t, err := tu.threadRep.FindByID(int32(ID))
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	return t, nil
}

func (tu *ThreadUsecase) VoteThread(slugID string, vote *models.Vote) (*models.Thread, error) {
	ID, err := strconv.Atoi(slugID)
	if err != nil {
		isExist := tu.threadRep.IsExistBySlug(slugID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this id doesnt't exist"
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}
		ID = tu.threadRep.GetIDBySlug(slugID)
	} else {
		isExist := tu.threadRep.IsExistByID(ID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this id doesnt't exist"
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}
	}

	isExistUser := tu.userUC.IsExistUser(vote.Nickname)
	if !isExistUser {
		rerr := new(models.HttpError)
		rerr.StringErr = "user doesn't exist"
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
	vote.Nickname = tu.userUC.GetRightNickname(vote.Nickname)

	isVotedYet := tu.threadRep.IsVoted(vote.Nickname, int32(ID))
	if !isVotedYet {
		err := tu.threadRep.CreateVote(int32(ID), vote)
		if err != nil {
			rerr := new(models.HttpError)
			rerr.StringErr = err.Error()
			rerr.StatusCode = http.StatusInternalServerError
			return nil, rerr
		}

		_ = tu.threadRep.AddVote(int32(ID), vote.Voice)

	} else {
		isVotedSame := tu.threadRep.IsVotedSame(vote, int32(ID))
		if !isVotedSame {
			err := tu.threadRep.EditVote(int32(ID), vote)
			if err != nil {
				rerr := new(models.HttpError)
				rerr.StringErr = err.Error()
				rerr.StatusCode = http.StatusInternalServerError
				return nil, rerr
			}
			_ = tu.threadRep.AddVote(int32(ID), vote.Voice * 2)
		}
	}

	 thr, err := tu.threadRep.FindByID(int32(ID))
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	return thr, nil
}

func (tu *ThreadUsecase) Status() int32 {
	return tu.threadRep.Status()
}

func (tu *ThreadUsecase) IsExistBySlug(slug string) bool {
	 return tu.threadRep.IsExistBySlug(slug)
}