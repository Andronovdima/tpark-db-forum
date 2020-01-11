package usecase

import (
	"github.com/Andronovdima/tpark-db-forum/internal/app/user/repository"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type UserUsecase struct {
	UserRep *repository.UserRepository
}

func NewUserUsecase(us *repository.UserRepository) *UserUsecase {
	UserUsecase := &UserUsecase{
		UserRep: us,
	}
	return UserUsecase
}

func (u *UserUsecase) CreateUser(user *models.User, nickname string) (*models.User, error) {
	err := new(models.HttpError)

	isExistUser := u.UserRep.IsExist(nickname)
	if isExistUser {
		err.StatusCode = http.StatusConflict
		err.StringErr = "user already exists with this email or nickname"

		us, cerr := u.UserRep.Find(nickname)
		if cerr != nil {
			err.StatusCode = http.StatusInternalServerError
			err.StringErr = "Internal error"
			return nil, err
		}

		return us, err
	}

	user.Nickname = nickname
	cerr := u.UserRep.Create(user)
	if cerr != nil {
		err.StatusCode = http.StatusInternalServerError
		err.StringErr = "Internal error"
	}

	return user, nil
}

func (u *UserUsecase) IsExistUser(nickname string) bool {
	return u.UserRep.IsExist(nickname)
}

func (u *UserUsecase) Find(nickname string) (*models.User, error) {
	us, cerr := u.UserRep.Find(nickname)
	if cerr != nil {
		return nil, cerr
	}
	return us, nil
}

func (u *UserUsecase) UpdateProfile(nickname string, user *models.User) (*models.User, error){
	rerr := new(models.HttpError)

	isExist := u.UserRep.IsExist(nickname)
	if !isExist {
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "Can't find user with that email and nickname"
		return nil, rerr
	}

	isBusyEmail := u.UserRep.CheckEmail(user.Email)
	if isBusyEmail {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "this email is exist yet"
		return nil, rerr
	}

	user.Nickname = nickname
	err := u.UserRep.Update(user)
	if  err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = "Internal Server Error"
		return nil, rerr
	}

	return user, nil
}
