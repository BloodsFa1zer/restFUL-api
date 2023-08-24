package hash

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

//func Verify(hashed, password string) bool {
//	user, err := cl.FindUser("Password", password)
//	if err != nil {
//		return false
//	}
//	errPassword := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(user.Password))
//	if errPassword == nil {
//		return true
//	}
//	return false
//}
