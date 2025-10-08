package password

import "golang.org/x/crypto/bcrypt"

const cost = 12

func Encrypt(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	return string(bytes), err
}

func Compare(hash, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
