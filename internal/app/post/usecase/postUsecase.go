package usecase

import (
	forumU "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	postR "github.com/Andronovdima/tpark-db-forum/internal/app/post/repository"
	threadU "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	userU "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strconv"
)

type PostUsecase struct {
	postRep  *postR.PostRepository
	threadUC *threadU.ThreadUsecase
	forumUC  *forumU.ForumUsecase
	userUC  *userU.UserUsecase
}

func NewPostUsecase(pr *postR.PostRepository, tr *threadU.ThreadUsecase, fr *forumU.ForumUsecase, ur *userU.UserUsecase) *PostUsecase {
	postUsecase := &PostUsecase{
		postRep:  pr,
		threadUC: tr,
		forumUC: fr,
		userUC: ur,
	}
	return postUsecase
}

func (p *PostUsecase) CreatePosts(threadSlugID string, posts *[]models.Post) (*[]models.Post, error) {
	var isExist bool
	threadID, err := strconv.Atoi(threadSlugID)
	if err == nil {
		isExist = p.threadUC.IsExistByID(threadID)
	} else {
		isExist = p.threadUC.IsExistBySlug(threadSlugID)
		if isExist {
			threadID = p.threadUC.GetIDBySlug(threadSlugID)
		}
	}

	if !isExist {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Thread with this slug doesn't exist"
		return nil, rerr
	}

	forum, err := p.threadUC.GetForum("", int32(threadID))
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = rerr.Error()
		return nil, rerr
	}

	for _, pr := range *posts {
		if pr.Parent == 0 {
			continue
		}
		isExist = p.postRep.IsExist(pr.Parent)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StatusCode = http.StatusConflict
			rerr.StringErr = "error with parrentID and ID"
			return nil, rerr
		}
	}

	code , err := p.postRep.CreatePosts(int32(threadID), forum, posts)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = code
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	err = p.forumUC.AddPost(forum, int64(len(*posts)))
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	return posts, nil
}

func (p *PostUsecase) EditPost(id int64, update *models.PostUpdate) (*models.Post, error) {
	isExist := p.postRep.IsExist(id)
	if !isExist {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "post with this id doesn't exist"
		return nil, rerr
	}

	updPost, err := p.postRep.Find(id)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	if update.Message != "" && update.Message != updPost.Message {
		err := p.postRep.EditPost(id, update)
		if err != nil {
			rerr := new(models.HttpError)
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, rerr
		}
	}

	updPost, err = p.postRep.Find(id)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	return updPost, nil
}

func (p *PostUsecase) GetInfo(id int64, includeForum bool, includeThread bool, includeUser bool) (*models.PostFull, error) {
	isExist := p.postRep.IsExist(id)
	if !isExist {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "post with this id doesn't exist"
		return nil, rerr
	}

	postFull := new(models.PostFull)

	currPost, err := p.postRep.Find(id)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	postFull.Post = currPost

	if includeForum {
		forum, err := p.forumUC.Find(currPost.Forum)
		if err != nil {
			rerr := new(models.HttpError)
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, rerr
		}
		postFull.Forum = forum
	}

	if includeThread {
		thread, err := p.threadUC.GetThreadByID(currPost.Thread)
		if err != nil {
			rerr := new(models.HttpError)
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, rerr
		}
		postFull.Thread = thread
	}

	if includeUser {
		user, err := p.userUC.Find(currPost.Author)
		if err != nil {
			rerr := new(models.HttpError)
			rerr.StatusCode = http.StatusInternalServerError
			rerr.StringErr = err.Error()
			return nil, rerr
		}
		postFull.Author = user
	}

	return postFull, nil
}

func (p *PostUsecase) GetPosts(slugID string, limit int, since int, sort string, desc bool) ([]*models.Post, error) {
	ID, err := strconv.Atoi(slugID)
	if err != nil {
		isExist := p.threadUC.IsExistBySlug(slugID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this slug doesnt't exist" + slugID
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}

		ID = p.threadUC.GetIDBySlug(slugID)
	} else {
		isExist := p.threadUC.IsExistByID(ID)
		if !isExist {
			rerr := new(models.HttpError)
			rerr.StringErr = "thread with this id doesnt't exist " + strconv.Itoa(ID)
			rerr.StatusCode = http.StatusNotFound
			return nil, rerr
		}
	}

	posts, err := p.postRep.GetThreadPosts(int32(ID), limit , since, sort, desc)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}

	return posts, nil
}


func (p *PostUsecase) Status() int64 {
	return p.postRep.Status()
}

//func (p PostUsecase) GetThreadPosts(slugID string, params map[string][]string) ([]*models.Post, error) {
//		ID, err := strconv.Atoi(slugID)
//		if err != nil {
//			isExist := p.threadUC.IsExistBySlug(slugID)
//			if !isExist {
//				rerr := new(models.HttpError)
//				rerr.StringErr = "thread with this slug doesnt't exist" + slugID
//				rerr.StatusCode = http.StatusNotFound
//				return nil, rerr
//			}
//
//			ID = p.threadUC.GetIDBySlug(slugID)
//		} else {
//			isExist := p.threadUC.IsExistByID(ID)
//			if !isExist {
//				rerr := new(models.HttpError)
//				rerr.StringErr = "thread with this id doesnt't exist " + strconv.Itoa(ID)
//				rerr.StatusCode = http.StatusNotFound
//				return nil, rerr
//			}
//		}
//	threadObj, err := p.threadUC.GetThreadByID(int32(ID))
//	if err != nil {
//		return nil, err
//	}
//
//	limit := "100"
//	if len(params["limit"]) >= 1 {
//		limit = params["limit"][0]
//	}
//	desc := ""
//	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
//		desc = "desc"
//	}
//	since := ""
//	if len(params["since"]) >= 1 {
//		since = params["since"][0]
//	}
//	sort := "flat"
//	if len(params["sort"]) >= 1 {
//		sort = params["sort"][0]
//	}
//
//	posts, err := p.postRep.GetThreadPosts(*threadObj, limit, desc, since, sort)
//	if err != nil {
//		return nil, err
//	}
//
//	return posts, nil
//}