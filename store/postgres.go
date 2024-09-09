package store

import (
	"database/sql"
	"document-service/models"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

const (
	sqlSelectDocument = `SELECT files.id, name, file, public, mime, files.create_date, json_agg(json_build_object('login', users.login)) FROM service_documents.files
    INNER JOIN service_documents.users_file_link ufl on files.id = ufl.id_document
    INNER JOIN service_documents.users ON users.id = ufl.id_user GROUP BY files.id`
)

// InitDB инициализирует соединение с базой данных
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName) // Создаем подключение
	if err != nil {
		return err
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func PostConn() (*sql.DB, error) {
	connect := "user=postgres password=1234567890 dbname=documents sslmode=disable"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

func CreateUser(user models.User) error {
	// Подключение к PostgreSQL
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)

	res, err := db.Prepare(`INSERT INTO service_documents.users(login, password) VALUES ($1, crypt($2, gen_salt('bf')))`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec(user.Login, user.Password)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func AuthenticateUser(login, password string) (bool, error) {
	// Аутентификация пользователя
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return false, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	var userExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT * FROM service_documents.users where login= $1 and password = crypt($2, password)) as user_exist`, login, password).Scan(&userExists)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return userExists, nil
}

// SaveDocument сохраняет документ в базу данных PostgreSQL
func SaveDocument(doc models.Document) (models.Document, error) {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return models.Document{}, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)

	query := `
    INSERT INTO service_documents.files (name, file, public, mime)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (id) DO NOTHING
    RETURNING id;
  `
	err = db.QueryRow(query, doc.Name, doc.File, doc.Public, doc.Mime).Scan(&doc.ID)
	if err != nil {
		log.Println(err)
		return models.Document{}, err
	}
	return doc, nil
}

func GetDocumentList(limit, offset uint64) (models.DocumentsList, error) {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return models.DocumentsList{}, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	var params string
	if limit != 0 {
		params = fmt.Sprintf(" limit %d", limit)
		if offset != 0 {
			params = fmt.Sprintf("%s offset %d", params, offset)
		}
	}
	rows, err := db.Query(sqlSelectDocument + params)
	if err != nil {
		log.Println(err)
		return models.DocumentsList{}, err
	}
	defer rows.Close()

	var documentsList models.DocumentsList

	for rows.Next() {
		var jsonData []byte
		var doc models.Document
		err := rows.Scan(&doc.ID, &doc.Name, &doc.File, &doc.Public, &doc.Mime, &doc.Created, &jsonData)
		if err != nil {
			log.Println(err)
			return models.DocumentsList{}, err
		}
		err = json.Unmarshal(jsonData, &doc.UsersFileLinkList)
		if err != nil {
			log.Println("Ошибка при парсинге JSON:", err)
		}
		documentsList = append(documentsList, doc)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return models.DocumentsList{}, err
	}

	return documentsList, nil
}

func GetDocument(id uint64) (models.Document, error) {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return models.Document{}, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	var jsonData []byte
	var doc models.Document
	err = db.QueryRow(sqlSelectDocument+` HAVING files.id = $1`, id).Scan(&doc.ID, &doc.Name, &doc.File, &doc.Public, &doc.Mime, &doc.Created, &jsonData)
	if err != nil {
		log.Println(err)
		return models.Document{}, err
	}
	err = json.Unmarshal(jsonData, &doc.UsersFileLinkList)
	if err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
	}
	return doc, nil
}

// Удаление документа
func DeleteDocument(id uint64) error {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	res, err := db.Prepare(`DELETE FROM service_documents.files WHERE files.id = $1`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetUserIDByLogin(login string) (uint64, error) {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return 0, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	var id uint64
	err = db.QueryRow(`SELECT id FROM users WHERE login = $1`, login).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}
func SaveUserFileLink(userFileLink models.UsersFileLink) error {
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	res, err := db.Prepare(`INSERT INTO service_documents.users_file_link(id_user, id_document) VALUES ($1, $2)`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec(userFileLink.IdUser, userFileLink.IdDocument)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
