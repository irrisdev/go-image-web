CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    subject TEXT,
    message TEXT,
    image_uuid TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
