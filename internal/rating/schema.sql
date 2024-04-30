CREATE TABLE IF NOT EXISTS memes
(
    id         TEXT    NOT NULL PRIMARY KEY,
    channel_id TEXT    NOT NULL,
    member_id  TEXT    NOT NULL,
    score      INTEGER NOT NULL,
    timestamp  FLOAT   NOT NULL,
    link       TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS members
(
    id           TEXT NOT NULL PRIMARY KEY,
    full_name    TEXT NOT NULL,
    display_name TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_channel_id ON memes (channel_id);
