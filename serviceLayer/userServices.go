package serviceLayer

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"time"
)

type UserService struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

func NewUserService(DbUser database.DbInterface, validate *validator.Validate) *UserService {
	return &UserService{DbUser: DbUser, validate: validate}
}

func (us *UserService) CreateUser(user database.User) (int64, error) {
	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	if err := us.UserValidation(newUser); err != nil {
		return 0, err
	}

	return us.DbUser.InsertUser(newUser)
}

func (us *UserService) GetUser(userID int64) (*database.User, error) {
	return us.DbUser.FindByID(userID)
}

func (us *UserService) EditUser(ID int64, user database.User) (int64, error) {
	if err := us.UserValidation(user); err != nil {
		return 0, err
	}

	return us.DbUser.UpdateUser(ID, user)
}

func (us *UserService) GetAllUsers() (*[]database.User, error) {
	return us.DbUser.FindUsers()
}

func (us *UserService) DeleteUser(userID int64) error {
	return us.DbUser.DeleteUserByID(userID)
}

func (us *UserService) GetUserByName(nickname, password string) (string, error) {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	user, err := us.DbUser.FindByNicknameToGetUserPassword(nickname)

	if hash.Verify(user.Password, password) != true {
		return "", echo.ErrUnauthorized
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

	return t, err
}

func (us *UserService) UserValidation(user database.User) error {

	if validationErr := us.validate.Struct(&user); validationErr != nil {
		return validationErr
	}
	return nil
}

func (us *UserService) IsUserHavePermission(roleToCheck string, user interface{}) bool {
	userToken, ok := user.(*jwt.Token)
	if !userToken.Valid {
		return false
	}
	if !ok {
		return false
	}
	claims := userToken.Claims.(*config.JwtCustomClaims)

	return claims.Role == roleToCheck
}
