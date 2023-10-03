PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Voting (
    user_id INTEGER DEFAULT 0,
    voter_id INTEGER DEFAULT 0,
    updated_at TEXT DEFAULT '0',
    vote_value INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES Users(ID),
    FOREIGN KEY (voter_id) REFERENCES Users(ID)
);

-- INSERT INTO Voting (user_id, voter_id, vote_value) VALUES (11,15, -1);

SELECT sum(vote_value) from Voting


