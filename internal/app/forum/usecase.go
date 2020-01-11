package forum

import "github.com/Andronovdima/tpark-db-forum/internal/models"

type Usecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	CreateThread(th *models.Thread, slug string) (*models.Thread, error)
}