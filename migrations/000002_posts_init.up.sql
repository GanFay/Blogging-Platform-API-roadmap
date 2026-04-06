CREATE TABLE IF NOT EXISTS posts (
                                     id         SERIAL PRIMARY KEY,
                                     author_id  INT	NOT NULL,
                                     title      TEXT UNIQUE NOT NULL,
                                     content    TEXT        NOT NULL,
                                     category   TEXT,
                                     tags TEXT[],
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);