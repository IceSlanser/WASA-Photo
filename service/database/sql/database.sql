-- Create profiles table
CREATE TABLE IF NOT EXISTS profiles (
    ID INTEGER PRIMARY KEY,
    Username TEXT UNIQUE,
    FollowingCount INTEGER DEFAULT 0,
    FollowerCount INTEGER DEFAULT 0,
    PostCount INTEGER DEFAULT 0
);

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    ID INTEGER PRIMARY KEY,
    ProfileID INTEGER NOT NULL,
    File BLOB,
    Description TEXT NOT NULL,
    LikeCount INTEGER DEFAULT 0,
    CommentCount INTEGER DEFAULT 0,
    DateTime DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (ProfileID) REFERENCES profiles(ID)
);

-- Create follows table
CREATE TABLE IF NOT EXISTS follows (
    ID INTEGER PRIMARY KEY,
    FollowerUID INTEGER,
    FollowedUID INTEGER,
    UNIQUE (FollowerUID, FollowedUID),
    FOREIGN KEY (FollowerUID) REFERENCES profiles(ID) ON DELETE CASCADE,
    FOREIGN KEY (FollowedUID) REFERENCES profiles(ID) ON DELETE CASCADE
);

-- Create bans table
CREATE TABLE IF NOT EXISTS bans (
    ID INTEGER PRIMARY KEY,
    BannerUID INTEGER,
    BannedUID INTEGER,
    UNIQUE (BannerUID, BannedUID),
    FOREIGN KEY (BannerUID) REFERENCES profiles(ID) ON DELETE CASCADE,
    FOREIGN KEY (BannedUID) REFERENCES profiles(ID) ON DELETE CASCADE
);

-- Create likes table
CREATE TABLE IF NOT EXISTS likes (
    ID INTEGER PRIMARY KEY,
    PostID INTEGER,
    OwnerID INTEGER,
    UNIQUE (PostID, OwnerID),
    FOREIGN KEY (PostID) REFERENCES posts(ID) ON DELETE CASCADE,
    FOREIGN KEY (OwnerID) REFERENCES profiles(ID) ON DELETE CASCADE
);

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    ID INTEGER PRIMARY KEY,
    PostID INTEGER NOT NULL,
    OwnerID INTEGER NOT NULL,
    Text TEXT NOT NULL,
    DateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (PostID) REFERENCES posts(ID) ON DELETE CASCADE,
    FOREIGN KEY (OwnerID) REFERENCES profiles(ID) ON DELETE CASCADE
);

-- Insert multiple values into the profiles table
INSERT INTO profiles (Username) VALUES
    ('Chino'),
    ('Pala'),
    ('Saidai'),
    ('Cavito'),
    ('Hermano');

-- Insert multiple values into the followers table
INSERT INTO follows (FollowerUID, FollowedUID) VALUES
    (5, 4),
    (4, 5),
    (4, 2);

INSERT INTO posts (ProfileID) VALUES
    (4), (4), (4), (4);

INSERT INTO bans (BannerUID, BannedUID) VALUES
        (3, 1);

INSERT INTO comments (PostID, OwnerID, Text) VALUES
    (1, 5, 'Soyeon MLML'),
    (1, 4, 'IL REEEEE'),
    (1, 3, 'Bomb on Kiev'),
    (1, 1, 'Kiryu enjoyer');

INSERT INTO likes (PostID, OwnerID) VALUES
    (1, 1), (1, 3), (1, 4), (1, 5), (1, 2);

-- Select
SELECT * FROM profiles;
SELECT * FROM follows;
SELECT * FROM posts;
SELECT * FROM bans;

-- Try
SELECT FollowedUID
FROM follows
WHERE FollowerUID = 3 AND FollowedUID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = 1)