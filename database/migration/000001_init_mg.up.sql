ALTER TABLE Users
ADD COLUMN role TEXT DEFAULT 'User';

CREATE TABLE IF NOT EXISTS Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT ,
    nick_name TEXT UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT 'NULL',
    deleted_at TEXT DEFAULT 'NULL'

)