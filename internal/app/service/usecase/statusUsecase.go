package usecase

import (
	forumU "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	threadU "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	userU "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
)

type ServiceUsecase struct {
	threadUC *threadU.ThreadUsecase
	forumUC  *forumU.ForumUsecase
	userUC  *userU.UserUsecase
}

func NewServiceUsecase(tr *threadU.ThreadUsecase, fr *forumU.ForumUsecase, ur *userU.UserUsecase) *ServiceUsecase {
	postUsecase := &ServiceUsecase{
		threadUC: tr,
		forumUC: fr,
		userUC: ur,
	}
	return postUsecase
}