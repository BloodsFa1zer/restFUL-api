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

var userVotes = make(map[string][]int64)
var votesTime = make(map[string]time.Time)

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

func (us *UserService) CreateToken(nickname, password string) (string, error, int) {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	user, err := us.DbUser.FindByNicknameToGetUserPassword(nickname)
	if err == sql.ErrNoRows {
		return "", errors.New("you have no account and will be redirected to registration page"), http.StatusSeeOther
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

func (us *UserService) Registration(username, firstName, surName, password string) (int, error, int) {
	user := database.User{
		Nickname:  username,
		FirstName: firstName,
		LastName:  surName,
		Password:  password,
	}
	if err := us.UserValidation(user); err != nil {
		return 0, err, http.StatusBadRequest
	}

	userID, err := us.DbUser.InsertUser(user)
	if err != nil {
		return 0, err, http.StatusInternalServerError
	}

	return int(userID), err, http.StatusCreated
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

func (us *UserService) GetUserNameViaToken(user interface{}) string {
	userToken, ok := user.(*jwt.Token)
	if !userToken.Valid {
		return ""
	}
	if !ok {
		return ""
	}
	claims := userToken.Claims.(*config.JwtCustomClaims)

	return claims.Name
}

func (us *UserService) Vote(userID int64, userName string) (error, int) {

	// TODO: here i need to check if user with that token does not vote for that person earlier and make a
	// time restriction for 1 hour

	userVote, voteTime, err := mapCreation(userName, userID)
	fmt.Println(voteTime)
	fmt.Println(userVote)

	//	if _, ok := voteTime[userName]; ok {
	// fmt.Println("eff")
	if !us.isUserAllowedToVoteAgain(voteTime[userName]) {
		return errors.New("user only allowed to vote once in an hour, your last vote was at:" + voteTime[userName].Format("2006.01.02 15:04")), http.StatusLocked
	}
	//}

	//if !us.isUserAllowedToVoteForThatCandidate(userVote, userName, userID) {
	//	return errors.New("you are not allowed to vote for the same candidate twice"), http.StatusLocked
	//}

	err = us.DbUser.VoteForUser(userID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return err, http.StatusOK
}

func (us *UserService) GetUserRate(ID int64) (*database.UserRating, error, int) {
	user, err := us.DbUser.GetUserRating(ID)
	if err == sql.ErrNoRows {
		return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return user, nil, http.StatusOK
}

// Maybe it is better to check it in the middleware?
func (us *UserService) isUserAllowedToVoteAgain(voteTime time.Time) bool {

	oneHourAgo := voteTime.Add(-1 * time.Hour)
	duration := voteTime.Sub(oneHourAgo)

	return duration > time.Hour
}

func isUserAllowedToVoteForThatCandidate(userVote map[string][]int64, userName string, desiredID int64) (bool, error) {
	votes := userVote[userName]
	fmt.Println("votes:", votes)
	for _, vote := range votes {
		if vote == desiredID {
			return false, errors.New("you are not allowed to vote for the same candidate twice")
		}
	}

	return true, nil
}

func mapCreation(userName string, voteID int64) (map[string][]int64, map[string]time.Time, error) {

	if _, ok := userVotes[userName]; ok {
		if isAllowed, err := isUserAllowedToVoteForThatCandidate(userVotes, userName, voteID); isAllowed {
			//userVotes[userName] = []int64{}

			userVotes[userName] = append(userVotes[userName], voteID)
			votesTime[userName] = time.Now()
		} else {
			// fmt.Println("you are not allowed to vote for the same candidate twice")
			return nil, votesTime, err
		}
	} else {
		userVotes[userName] = append(userVotes[userName], voteID)
	}
	return userVotes, votesTime, nil
}
