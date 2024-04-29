CREATE TABLE IF NOT EXISTS memes
(
    id         TEXT    NOT NULL PRIMARY KEY,
    channel_id TEXT    NOT NULL,
    member_id  TEXT    NOT NULL,
    score      INTEGER NOT NULL,
    timestamp  FLOAT   NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_channel_id ON memes (channel_id);
