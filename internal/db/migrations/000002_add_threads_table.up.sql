CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    uuid TEXT UNIQUE NOT NULL,
    author TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,

    board_id INTEGER,
    FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE SET NULL
);

CREATE INDEX threads_board_id_idx ON threads(board_id);