package handlers

import (
	"crypto/rand"
	"document-service/config"
	"document-service/models"
	"document-service/services"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// Регистрация нового пользователя
func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
		Login string `json:"login"`
		Pswd  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, nil, errors.New("invalid request"))
		return
	}

	// Проверка токена администратора
	if req.Token != config.Config.AdminToken {
		respondWithJSON(w, http.StatusUnauthorized, nil, errors.New("invalid token"))
		return
	}

	// Сохранение пользователя
	err := services.CreateUser(models.User{
		Login:    req.Login,
		Password: req.Pswd,
	})
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, nil, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"login": req.Login,
	}, nil)
}

// Аутентификация
func Auth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login string `json:"login"`
		Pswd  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, nil, errors.New("invalid request"))
		return
	}

	// Аутентификация пользователя
	token, err := services.AuthenticateUser(req.Login, req.Pswd)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, nil, err)
		return
	}

	services.Session[token] = struct{}{}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
	}, nil)
}

// Завершение сессии
func Logout(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	delete(services.Session, token)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		token: true,
	}, nil)
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := make(map[string]interface{})
	if err != nil {
		response["error"] = map[string]interface{}{
			"code": 500,
			"text": err.Error(),
		}
	}
	if data != nil {
		response["data"] = data
	}

	json.NewEncoder(w).Encode(response)
}
