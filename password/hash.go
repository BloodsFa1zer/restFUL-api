package password

import "golang.org/x/crypto/bcrypt"

func Hash(password []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, 14)
	return string(bytes), err
}
