CREATE TABLE bookmarks (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url        TEXT NOT NULL,
    title      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
