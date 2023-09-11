package services

import (
	"app3.1/database"
)

func (us *UserService) Create(user database.User) (int64, error) {
	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	return us.DbUser.InsertUser(newUser)
}

func (us *UserService) Get(userID int64) (*database.User, error) {
	return us.DbUser.FindByID(userID)
}

func (us *UserService) Edit(ID int64, user database.User) (int64, error) {
	return us.DbUser.UpdateUser(ID, user)
}

func (us *UserService) GetAll() (*[]database.User, error) {
	return us.DbUser.FindUsers()
}

func (us *UserService) Delete(userID int64) error {
	return us.DbUser.DeleteUserByID(userID)
}

func (us *UserService) GetPasswordByName(nickname string) (*database.User, error) {
	user, err := us.DbUser.FindByNicknameToGetUserPassword(nickname)
	return user, err
}

func (us *UserService) UserValidation(user database.User) error {

	if validationErr := us.validate.Struct(&user); validationErr != nil {
		return validationErr
	}
	return nil
}
