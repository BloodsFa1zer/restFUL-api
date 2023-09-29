package service

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"database/sql"
	"errors"
	"fmt"
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

func (us *UserService) CreateToken(user database.User) (string, error, int) {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	SelectedUser, err := us.DbUser.FindByNicknameToGetUserPassword(user.Nickname)
	if err == sql.ErrNoRows {
		return "", errors.New("you have no account and will be redirected to registration page"), http.StatusSeeOther
	} else if err != nil {
		return "", err, http.StatusInternalServerError
	}

	if hash.Verify(SelectedUser.Password, user.Password) != true {
		return "", errors.New("incorrect password"), http.StatusUnauthorized
	}

	claims := &config.JwtCustomClaims{
		ID:   SelectedUser.ID,
		Name: SelectedUser.Nickname,
		Role: SelectedUser.Role,
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

func (us *UserService) GetUserIDViaToken(user interface{}) (int64, error) {
	userToken, ok := user.(*jwt.Token)
	if !userToken.Valid {
		return 0, errors.New("token is invalid")
	}
	if !ok {
		return 0, errors.New("no such token found")
	}

	claims := userToken.Claims.(*config.JwtCustomClaims)
	fmt.Println(claims.Name)
	fmt.Println(claims.ID)

	return claims.ID, nil
}

func (us *UserService) PostVote(userID, voterID int) (error, int) {
	err := us.isUserAllowedToVote(voterID, userID)
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.WriteUserVotes(userID, voterID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) DeleteVote(userID, voterID int) (error, int) {
	err := us.isUserAllowedToVote(voterID, userID)
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.WriteUserVotes(userID, voterID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) isUserAllowedToVote(voterID, userID int) error {
	isAllowedCandidate := false

	voteTime, err := us.DbUser.GetUserVotes(int64(voterID), int64(userID))
	if err == sql.ErrNoRows {
		isAllowedCandidate = true
		voteTime = "0"
	} else if err != nil {
		return err
	}
	fmt.Println("vt:", voteTime)
	if voteTime == "0" {
		return nil
	}

	timeWhenUserVotes, err := castUserDataToUseInMap(voteTime)
	if err != nil {
		return err
	}

	if voterID == userID {
		return errors.New(" you are not allowed to vote for yourself")
	}

	if isAllowedCandidate {
		if !timeWhenUserVotes.IsZero() {
			errVoteTime := isUserAllowedToVoteAgainAfterOneHourTime(timeWhenUserVotes)
			if errVoteTime != nil {
				return errVoteTime
			}
		}
	} else {
		return errors.New("user cannot vote for the same candidate twice")
	}

	return nil
}

func isUserAllowedToVoteAgainAfterOneHourTime(timeWhenUserVotes time.Time) error {
	duration := time.Now().Sub(timeWhenUserVotes)
	if duration <= time.Hour {
		return errors.New("user only allowed to vote once in an hour, your last vote was at:" + timeWhenUserVotes.Format("2006.01.02 15:04"))
	}

	return nil
}

func castUserDataToUseInMap(timeWhenUserVote string) (time.Time, error) {

	voteTime := time.Time{}

	if timeWhenUserVote == "0" {
		return time.Time{}, nil
	}
	layout := "2006-01-02 15:04:05.000000-07:00"
	voteTime, err := time.Parse(layout, timeWhenUserVote)
	if err != nil {
		return time.Time{}, err
	}

	return voteTime, nil
}
