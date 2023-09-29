PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Voting (
    user_id INTEGER,
    voter_id INTEGER,
    updated_at TEXT DEFAULT '0',
    FOREIGN KEY (user_id) REFERENCES Users(ID),
    FOREIGN KEY (voter_id) REFERENCES Users(ID)
);

SELECT * FROM Voting WHERE user_id = 11;
