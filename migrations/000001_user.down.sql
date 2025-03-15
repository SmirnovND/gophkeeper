-- users: хранение данных о пользователях
DROP TABLE IF EXISTS users;

-- secrets: хранение приватных данных пользователей
DROP TABLE IF EXISTS secrets;

-- devices: хранение информации об авторизованных клиентах
DROP TABLE IF EXISTS devices;

-- audit_log: логирование действий пользователей
DROP TABLE IF EXISTS audit_log;