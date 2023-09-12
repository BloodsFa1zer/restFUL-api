package service

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type UserService struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

func NewUserService(DbUser database.DbInterface, validate *validator.Validate) *UserService {
	return &UserService{DbUser: DbUser, validate: validate}
}

func (us *UserService) CreateUser(user database.User) (int64, error, int) {
	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	if err := us.UserValidation(newUser); err != nil {
		return 0, err, http.StatusBadRequest
	}

	insertedId, err := us.DbUser.InsertUser(newUser)
	if err != nil {
		return insertedId, err, http.StatusInternalServerError
	}

	return insertedId, err, http.StatusCreated
}

func (us *UserService) GetUser(userID int64) (*database.User, error, int) {
	user, err := us.DbUser.FindByID(userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return user, nil, http.StatusOK
}

func (us *UserService) EditUser(ID int64, user database.User) (int64, error, int) {
	if err := us.UserValidation(user); err != nil {
		return 0, err, http.StatusBadRequest
	}
	updatedID, err := us.DbUser.UpdateUser(ID, user)
	if err == sql.ErrNoRows {
		return 0, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return 0, err, http.StatusInternalServerError
	}

	return updatedID, err, http.StatusOK
}

func (us *UserService) GetAllUsers() (*[]database.User, error, int) {

	users, err := us.DbUser.FindUsers()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return users, err, http.StatusOK
}

func (us *UserService) DeleteUser(userID int64) (error, int) {
	err := us.DbUser.DeleteUserByID(userID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return errors.New("user successfully deleted"), http.StatusOK
}

func (us *UserService) GetToken(nickname, password string) (string, error, int) {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	user, err := us.DbUser.FindByNicknameToGetUserPassword(nickname)
	if err == sql.ErrNoRows {
		return "", errors.New("there is no user with that NickName"), http.StatusBadRequest
	} else if err != nil {
		return "", err, http.StatusInternalServerError
	}

	if hash.Verify(user.Password, password) != true {
		return "", errors.New("incorrect password"), http.StatusUnauthorized
	}

	claims := &config.JwtCustomClaims{
		Name: user.Nickname,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(cfg.SigningKey))

	return t, err, http.StatusOK
}

func (us *UserService) UserValidation(user database.User) error {

	if validationErr := us.validate.Struct(&user); validationErr != nil {
		return validationErr
	}
	return nil
}

func (us *UserService) IsUserHavePermission(roleToCheck string, user interface{}) (bool, int) {
	userToken, ok := user.(*jwt.Token)
	if !userToken.Valid {
		return false, http.StatusBadRequest
	}
	if !ok {
		return false, http.StatusBadRequest
	}
	claims := userToken.Claims.(*config.JwtCustomClaims)

	return claims.Role == roleToCheck, http.StatusOK
}
