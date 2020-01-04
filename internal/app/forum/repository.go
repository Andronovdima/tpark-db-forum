package forum

import "github.com/Andronovdima/tpark-db-forum/internal/models"

type Repository interface {
	Create(forum *models.Forum) (*models.Forum, error)
}