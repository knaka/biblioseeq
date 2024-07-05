CREATE TABLE logs (
  id integer PRIMARY KEY AUTOINCREMENT,
  message text NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp
);
