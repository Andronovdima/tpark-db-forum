package usecase

import (
	forum "github.com/Andronovdima/tpark-db-forum/internal/app/forum/repository"
	"github.com/Andronovdima/tpark-db-forum/internal/app/user/repository"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type UserUsecase struct {
	UserRep *repository.UserRepository
	forumUC *forum.ForumRepository
}

func NewUserUsecase(us *repository.UserRepository, fr *forum.ForumRepository) *UserUsecase {
	UserUsecase := &UserUsecase{
		UserRep: us,
		forumUC: fr,
	}
	return UserUsecase
}

func (u *UserUsecase) CreateUser(user *models.User, nickname string) (*[]models.User, error) {
	err := new(models.HttpError)
	users := new([]models.User)

	isExistUserNickname := u.UserRep.IsExist(nickname)
	if isExistUserNickname {
		us, cerr := u.UserRep.Find(nickname)
		if cerr != nil {
			err.StatusCode = http.StatusInternalServerError
			err.StringErr = cerr.Error()
			return nil, err
		}

		*users = append(*users, *us)
	}

	isExistUserEmail := u.UserRep.IsExistEmail(user.Email)
	if isExistUserEmail {
		usr, cerr := u.UserRep.FindByEmail(user.Email)
		if cerr != nil {
			err.StatusCode = http.StatusInternalServerError
			err.StringErr = cerr.Error()
			return nil, err
		}
		if usr.Nickname != nickname {
			*users = append(*users, *usr)
		}
	}

	if isExistUserNickname || isExistUserEmail {
		err.StatusCode = http.StatusConflict
		err.StringErr = "user already exists with this email or nickname"
		return users, err
	}

	user.Nickname = nickname
	cerr := u.UserRep.Create(user)
	if cerr != nil {
		err.StatusCode = http.StatusInternalServerError
		err.StringErr = cerr.Error()
		return nil, err
	}
	*users = append(*users, *user)

	return users, nil
}

func (u *UserUsecase) IsExistUser(nickname string) bool {
	return u.UserRep.IsExist(nickname)
}

func (u *UserUsecase) Find(nickname string) (*models.User, error) {
	if !u.UserRep.IsExist(nickname) {
		rerr := new(models.HttpError)
		rerr.StringErr = "user with this nicknme doesnt exist"
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
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

	isBusyEmail := u.UserRep.IsExistEmail(user.Email)
	if isBusyEmail && user.Email != "" {
		rerr.StatusCode = http.StatusConflict
		rerr.StringErr = "this email is exist yet"
		return nil, rerr
	}

	user.Nickname = nickname
	us, err := u.UserRep.Find(nickname)
	if  err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = "Internal Server Error"
		return nil, rerr
	}

	if user.Email == "" {
		user.Email = us.Email
	}

	if user.About == "" {
		user.About = us.About
	}

	if user.Fullname == "" {
		user.Fullname = us.Fullname
	}

	err = u.UserRep.Update(user)
	if  err != nil {
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = "Internal Server Error"
		return nil, rerr
	}

	return user, nil
}

func (u *UserUsecase) GetForumUsers(slug string, limit int, since string, desc bool) ([]models.User, error){
	isExist := u.forumUC.IsExist(slug)
	if !isExist {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusNotFound
		rerr.StringErr = "forum with this slug doesn't exist"
		return nil, rerr
	}

	users, err := u.UserRep.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StatusCode = http.StatusInternalServerError
		rerr.StringErr = err.Error()
		return nil, rerr
	}

	return users, nil
}

func (u *UserUsecase) Status() int32 {
	return u.UserRep.Status()
}

func (u *UserUsecase) GetRightNickname(nickname string) string {
	return u.UserRep.GetRightNickname(nickname)
}

