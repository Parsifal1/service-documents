CREATE SCHEMA IF NOT EXISTS service_documents;

-- Создадим таблицу пользователей
CREATE TABLE IF NOT EXISTS service_documents.users
(
    "id"          SERIAL PRIMARY KEY,
    "login"       TEXT NOT NULL UNIQUE,
    "password"    TEXT NOT NULL,
    "create_date" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS service_documents.files
(
    "id"          SERIAL PRIMARY KEY,                 -- Уникальный идентификатор файла
    "name"        TEXT    NOT NULL,                   -- Имя файла
    "file"        BOOLEAN NOT NULL,                   -- Поле, указывающее, является ли запись файлом
    "public"      BOOLEAN NOT NULL,                   -- Поле, указывающее, доступен ли файл публично
    "mime"        TEXT    NOT NULL,                   -- MIME-тип файла
    "create_date" TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время создания записи
);
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS service_documents.users_file_link
(
    "id_user"     INTEGER REFERENCES service_documents.users on delete cascade,
    "id_document" INTEGER REFERENCES service_documents.files on delete cascade ,
    primary key (id_user, id_document)
)
