package handlers

import (
	"document-service/models"
	"document-service/services"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// Загрузка документа
func UploadDocument(w http.ResponseWriter, r *http.Request) {
	// Подтверждение авторизации
	token := r.Header.Get("Authorization")
	_, ok := services.Session[token]
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, nil, errors.New("unauthorized"))
		return
	}

	// Проверка и обработка формы
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, nil, err)
		return
	}

	// Загрузка файла
	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, nil, err)
		return
	}
	defer file.Close()
	buf := make([]byte, 8)
	if _, err := io.ReadFull(file, buf); err != nil {
		respondWithJSON(w, http.StatusBadRequest, nil, err)
		return
	}

	// Сохранение документа
	doc := models.Document{
		Name:     r.FormValue("name"),
		Mime:     r.FormValue("mime"),
		Token:    token,
		File:     true,
		FileName: header.Filename,
		Data:     buf,
		// Параметры grant и другие данные
	}
	doc, err = services.SaveDocument(doc)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, nil, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"file": header.Filename,
	}, nil)
}

// Получение списка документов
func GetListDocuments(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, ok := services.Session[token]
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, nil, errors.New("unauthorized"))
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	// Получение документов
	docs, err := services.GetListDocuments(uint64(limit), uint64(offset))
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, nil, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"docs": docs,
	}, nil)
}

// Получение одного документа
func GetDocument(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, ok := services.Session[token]
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, nil, errors.New("unauthorized"))
		return
	}

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	// Получение документа
	doc, err := services.GetDocument(uint64(id))
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, nil, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": doc,
	}, nil)
}

// Удаление документа
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, ok := services.Session[token]
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, nil, errors.New("unauthorized"))
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, nil, err)
	}

	// Удаление документа
	err = services.DeleteDocument(uint64(id))
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, nil, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{}, nil)
}
