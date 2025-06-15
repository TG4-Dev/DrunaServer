package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePassword(passwor1, password2 string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwor1),
		[]byte(password2),
	); err != nil {
		return err
	}
	return nil
}
