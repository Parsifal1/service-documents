package models

import (
	"fmt"
	"os"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserList []User

type Document struct {
	ID                uint64            `json:"id"`
	Name              string            `json:"name"`
	Mime              string            `json:"mime"`
	Token             string            `json:"token"`
	File              bool              `json:"file"`
	Public            bool              `json:"public"`
	Created           time.Time         `json:"created"`
	Grant             []string          `json:"grant"`
	FileName          string            `json:"_"`
	UsersFileLinkList UsersFileLinkList `json:"-"`
	Data              []byte            `json:"-"`
}
type DocumentsList []Document

func (doc Document) SaveFile() error {
	file, err := os.Create(doc.FileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return err
	}
	defer file.Close()

	// Запись данных в файл
	_, err = file.Write(doc.Data)
	if err != nil {
		fmt.Println("Ошибка при записи данных в файл:", err)
		return err
	}
	return nil
}

type UsersFileLink struct {
	IdUser     uint64 `json:"-"`
	IdDocument uint64 `json:"-"`
	Login      string `json:"login"`
}
type UsersFileLinkList []UsersFileLink
