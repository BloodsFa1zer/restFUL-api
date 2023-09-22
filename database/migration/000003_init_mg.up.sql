PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Voting (
    ID integer,
    votes []integer,
    rating integer,
    vote_time text,
    FOREIGN KEY (ID) REFERENCES Users(ID)
)
