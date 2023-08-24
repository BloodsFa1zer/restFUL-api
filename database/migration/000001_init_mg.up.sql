CREATE TABLE IF NOT EXISTS Users (
    ID integer primary key autoincrement,
    NickName   text UNIQUE,
    FirstName TEXT not null,
    LastName TEXT not null,
    Password text not null,
    created_at text not null,
    updated_at text DEFAULT 'Unknown',
    deleted_at text DEFAULT 'Unknown'
)