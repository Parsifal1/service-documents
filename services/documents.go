package services

import (
	"document-service/models"
	"document-service/store"
	"errors"
)

var Documents map[uint64]models.Document

// SaveDocument сохраняет документ в памяти (вместо базы данных)
func SaveDocument(doc models.Document) (models.Document, error) {
	// Здесь может быть логика для записи в базу данных или другое хранилище
	_, ok := Session[doc.Token]
	if !ok {
		return models.Document{}, errors.New("no token in document")
	}
	document, err := store.SaveDocument(doc)
	if err != nil {
		return document, err
	}
	for _, item := range doc.Grant {
		//Выбрать id пользователя по логину
		userID, err := store.GetUserIDByLogin(item)
		if err != nil {
			return document, err
		}
		//Записать в userFileLink id пользователя и document.id
		err = store.SaveUserFileLink(models.UsersFileLink{
			IdUser:     userID,
			IdDocument: document.ID,
		})
		if err != nil {
			return document, err
		}
	}
	//Сохранение на диск
	err = document.SaveFile()
	if err != nil {
		return document, err
	}
	return document, nil
}

// GetListDocuments возвращает список документов по фильтрам
func GetListDocuments(limit, offset uint64) ([]models.Document, error) {
	// Если лимит указан, вернуть соответствующее количество документов
	// Если offset, то отдаем порцию со смещением
	var result []models.Document
	result, err := store.GetDocumentList(limit, offset)
	if err != nil {
		return result, err
	}
	for i, document := range result {
		for _, grant := range document.UsersFileLinkList {
			result[i].Grant = append(result[i].Grant, grant.Login)
		}
	}
	return result, nil
}

// GetDocument возвращает документ по ID
func GetDocument(id uint64) (models.Document, error) {
	doc, ok := Documents[id]
	if ok {
		return doc, nil
	}
	document, err := store.GetDocument(id)
	if err != nil {
		return document, err
	}
	for _, grant := range document.UsersFileLinkList {
		document.Grant = append(document.Grant, grant.Login)
	}
	Documents[id] = document
	return document, nil
}

// DeleteDocument удаляет документ по ID
func DeleteDocument(id uint64) error {
	err := store.DeleteDocument(id)
	if err != nil {
		return err
	}
	delete(Documents, id)
	return nil
}
