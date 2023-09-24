PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Voting (
    userID integer primary key,
    user_nick text,
    votes text DEFAULT '0',
    rating integer DEFAULT 0,
    vote_time text DEFAULT '0',
    FOREIGN KEY (userID) REFERENCES Users(ID)
);

INSERT INTO Voting (userID, user_nick)
SELECT ID, nick_name FROM Users
WHERE deleted_at = 'NULL';
