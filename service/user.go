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
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

var userVotes = make(map[string][]int)
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

func (us *UserService) GetUserNameViaToken(user interface{}) string {
	userToken, ok := user.(*jwt.Token)
	if !userToken.Valid {
		return ""
	}
	if !ok {
		return ""
	}
	claims := userToken.Claims.(*config.JwtCustomClaims)
	fmt.Println(claims.Name)

	return claims.Name
}

func (us *UserService) PostVote(userID int64, userName string) (error, int) {
	userVote, voteTime, err := us.isUserAllowedToVoteFor(userName, int(userID))
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.VoteForUser(userID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	err = us.DbUser.WriteUserVotes(voteTime, userVote, userName)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) DeleteVote(userID int64, userName string) (error, int) {
	userVote, voteTime, err := us.isUserAllowedToVoteAgainst(userName, int(userID))
	if err != nil {
		return err, http.StatusLocked
	}

	err = us.DbUser.VoteAgainstUser(userID)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	err = us.DbUser.WriteUserVotes(voteTime, userVote, userName)
	if err == sql.ErrNoRows {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (us *UserService) GetUserRateModerator(ID int64) (*database.UserRating, error, int) {
	userRating, err := us.DbUser.GetModeratorUserRating(ID)
	if err == sql.ErrNoRows {
		return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return userRating, nil, http.StatusOK
}

func (us *UserService) GetUserRate(ID int64) (*database.UserRating, error, int) {
	userRating, err := us.DbUser.GetUserRating(ID)
	if err == sql.ErrNoRows {
		return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return userRating, nil, http.StatusOK
}

func (us *UserService) GetAllUsersRate() (*[]database.UserRating, error, int) {
	userRating, err := us.DbUser.GetAllUsersRating()
	if err == sql.ErrNoRows {
		return nil, errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return userRating, nil, http.StatusOK
}

func (us *UserService) isUserAllowedToVoteFor(userName string, voteID int) (map[string][]int, map[string]time.Time, error) {
	votesTime[userName] = time.Now()
	userVotes[userName] = []int{}

	user, err := us.DbUser.GetUserVotes(userName)
	if err != nil {
		return nil, nil, err
	}

	votesOfUser, timeWhenUserVotes, err := castUserDataToUseInMap(user.UserVotes, user.VoteTime)
	if err != nil {
		return nil, nil, err
	}

	if int(user.ID) == voteID {
		return nil, nil, errors.New(" you are not allowed to vote for yourself")
	}

	isUserAllowedTime := true
	errVoteAgain := error(nil)
	if !timeWhenUserVotes.IsZero() {
		isUserAllowedTime, errVoteAgain = isUserAllowedToVoteAgain(timeWhenUserVotes)
	}
	timeWhenUserVotes = time.Now()

	isAllowed, errVoteTime := isUserAllowedToVoteForThatCandidate(votesOfUser, voteID)

	if isUserAllowedTime && isAllowed {
		for _, num := range votesOfUser {
			userVotes[userName] = append(userVotes[userName], num)
		}
		userVotes[userName] = append(userVotes[userName], voteID)
		votesTime[userName] = timeWhenUserVotes
	} else {
		if errVoteAgain != nil {
			return nil, votesTime, errVoteAgain
		} else {
			return nil, votesTime, errVoteTime
		}
	}

	fmt.Println("votesFromMap:", userVotes[userName])
	fmt.Println("time:", votesTime[userName])

	return userVotes, votesTime, nil
}

func (us *UserService) isUserAllowedToVoteAgainst(userName string, voteID int) (map[string][]int, map[string]time.Time, error) {
	votesTime[userName] = time.Now()
	userVotes[userName] = []int{}

	user, err := us.DbUser.GetUserVotes(userName)
	if err != nil {
		return nil, nil, err
	}

	votesOfUser, timeWhenUserVotes, err := castUserDataToUseInMap(user.UserVotes, user.VoteTime)
	if err != nil {
		return nil, nil, err
	}

	isUserAllowed := true
	errVoteAgain := error(nil)
	if !timeWhenUserVotes.IsZero() {
		isUserAllowed, errVoteAgain = isUserAllowedToVoteAgain(timeWhenUserVotes)
	}
	timeWhenUserVotes = time.Now()

	if isUserAllowed {
		for _, num := range votesOfUser {
			if voteID != num {
				userVotes[userName] = append(userVotes[userName], num)
			}
		}
		votesTime[userName] = timeWhenUserVotes
	} else {
		if errVoteAgain != nil {
			return nil, votesTime, errVoteAgain
		}
	}

	fmt.Println("votesFromMap:", userVotes[userName])
	fmt.Println("time:", votesTime[userName])

	return userVotes, votesTime, nil
}

func isUserAllowedToVoteAgain(timeWhenUserVotes time.Time) (bool, error) {
	duration := time.Now().Sub(timeWhenUserVotes)
	if duration <= time.Hour {
		return false, errors.New("user only allowed to vote once in an hour, your last vote was at:" + timeWhenUserVotes.Format("2006.01.02 15:04"))
	}

	return true, nil
}

func isUserAllowedToVoteForThatCandidate(votesOfUser []int, desiredID int) (bool, error) {
	for _, vote := range votesOfUser {
		if vote == desiredID {
			return false, errors.New("you are not allowed to vote for the same candidate twice")
		}
	}

	return true, nil
}

func castUserDataToUseInMap(userVotes, timeWhenUserVote string) ([]int, time.Time, error) {

	voteTime := time.Time{}
	var votes []int
	numStrParts := strings.Split(userVotes, ", ")
	for _, numStrPart := range numStrParts {
		num, err := strconv.Atoi(numStrPart)
		if err != nil {
			return nil, time.Time{}, err
		}
		votes = append(votes, num)
	}

	if timeWhenUserVote == "0" {
		return votes, time.Time{}, nil
	}
	layout := "2006-01-02 15:04:05.000000-07:00"
	voteTime, err := time.Parse(layout, timeWhenUserVote)
	if err != nil {
		return nil, time.Time{}, err
	}

	return votes, voteTime, nil
}
