package handlers

import (
	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) {
	// Регистрация пользователя
	router.HandleFunc("/api/register", Register).Methods("POST")

	// Аутентификация
	router.HandleFunc("/api/auth", Auth).Methods("POST")

	// Загрузка документа
	router.HandleFunc("/api/docs", UploadDocument).Methods("POST")

	// Список документов
	router.HandleFunc("/api/docs", GetListDocuments).Methods("GET", "HEAD")

	// Получение одного документа
	router.HandleFunc("/api/docs/{id}", GetDocument).Methods("GET", "HEAD")

	// Удаление документа
	router.HandleFunc("/api/docs/{id}", DeleteDocument).Methods("DELETE")

	// Завершение сессии
	router.HandleFunc("/api/auth/{token}", Logout).Methods("DELETE")
}
