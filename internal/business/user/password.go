package user

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CheckPassword(loginPass string, registeredPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(registeredPass), []byte(loginPass))
	if err != nil {
		return false
	}
	return true
}