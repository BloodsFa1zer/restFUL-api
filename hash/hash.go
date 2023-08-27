package hash

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	if err != nil {
		log.Panic().Err(err).Msg("can`t hash user password")
	}
	return string(bytes)
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
