package services

import (
	"document-service/models"
	"document-service/store"
	"errors"
	"math/rand"
	"regexp"
	"unicode/utf8"
)

var Session map[string]struct{}

const (
	keyLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func CreateUser(user models.User) error {
	// Проверка логина и пароля
	if !validateLogin(user.Login) {
		return errors.New("invalid login")
	}
	if !validatePassword(user.Password) {
		return errors.New("invalid password")
	}
	// Создание нового пользователя
	err := store.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func genToken() string {
	bytesSlice := make([]byte, 16)
	for i := range bytesSlice {
		bytesSlice[i] = keyLetter[rand.Int63()%int64(len(keyLetter))]
	}
	return string(bytesSlice)
}

func AuthenticateUser(login, password string) (string, error) {
	// Аутентификация пользователя
	res, err := store.AuthenticateUser(login, password)
	if err != nil {
		return "", err
	}
	if !res {
		return "", errors.New("неверно указан логин или пароль")
	}
	return genToken(), nil
}

func validateLogin(login string) bool {
	if utf8.RuneCountInString(login) < 8 {
		return false
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	return regex.MatchString(login)
}

func validatePassword(login string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	return regex.MatchString(login)
}
