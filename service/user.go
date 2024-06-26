package service

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"app3.1/redisDB"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type UserService struct {
	DbUser      database.DbInterface
	validate    *validator.Validate
	RedisClient redisDB.ClientRedisInterface
}

func NewUserService(DbUser database.DbInterface, validate *validator.Validate, RedisClient redisDB.ClientRedisInterface) *UserService {
	return &UserService{DbUser: DbUser, validate: validate, RedisClient: RedisClient}
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
	user, err := us.RedisClient.GetUser(userID)

	if err == redis.Nil {
		user, err = us.DbUser.FindByID(userID)
		if err == sql.ErrNoRows {
			return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
		} else if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		err = us.RedisClient.SetUser(*user)
		if err != nil {
			return nil, err, http.StatusLocked
		}
	} else if err != nil {
		return nil, err, http.StatusBadRequest
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
	users, err := us.RedisClient.GetUsers()
	if err == redis.Nil {
		users, err = us.DbUser.FindUsers()
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		err = us.RedisClient.SetUsers(*users)
		if err != nil {
			return nil, err, http.StatusLocked
		}
	} else if err != nil {
		return nil, err, http.StatusBadRequest
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

func (us *UserService) PostVoteFor(userID, voterID int) (error, int) {
	_, err := us.isUserAllowedToVote(voterID, userID)
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.WriteUserVotes(userID, voterID, 1)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) PostVoteAgainst(userID, voterID int) (error, int) {
	_, err := us.isUserAllowedToVote(voterID, userID)
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.WriteUserVotes(userID, voterID, -1)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) DeleteVote(userID, voterID int) (error, int) {
	err := us.DbUser.WithdrawVote(userID, voterID)
	if err == sql.ErrNoRows {
		return errors.New("there is no such vote to delete"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) ChangeVote(userID, voterID int) (error, int) {
	err := us.DbUser.ChangeVote(userID, voterID)
	if err == sql.ErrNoRows {
		return errors.New("there is no such vote to delete"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) isUserAllowedToVote(voterID, userID int) (bool, error) {
	if voterID == userID {
		return false, errors.New(" you are not allowed to vote for yourself")
	}

	exists, err := us.DbUser.IsSuchVoteExists(userID, voterID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if exists {
		return false, errors.New("user cannot vote for the same candidate twice")
	}

	voteTime, err := us.DbUser.GetUserLastVoteTime(voterID)
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return false, err
	}

	timeWhenUserVotes, err := castUserTime(voteTime)
	if err != nil {
		return false, err
	}

	if !timeWhenUserVotes.IsZero() {
		duration := time.Now().Sub(timeWhenUserVotes)
		if duration <= time.Hour {
			return false, errors.New("user only allowed to vote once in an hour, your last vote was at:" + timeWhenUserVotes.Format("2006.01.02 15:04"))
		}
	}

	return true, nil
}

func castUserTime(timeWhenUserVote string) (time.Time, error) {

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
