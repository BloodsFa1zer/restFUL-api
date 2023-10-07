PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Voting (
    user_id INTEGER DEFAULT 0,
    voter_id INTEGER DEFAULT 0,
    updated_at TEXT DEFAULT '0',
    vote_value INTEGER DEFAULT 0 CHECK (vote_value IN (-1, 1)),
    FOREIGN KEY (user_id) REFERENCES Users(ID),
    FOREIGN KEY (voter_id) REFERENCES Users(ID)
);


SELECT Users.*, sum(Voting.vote_value) from Users
left join Voting on Voting.user_id = Users.ID where ID = ?;




