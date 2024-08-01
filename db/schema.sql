CREATE TABLE logs (
  id integer PRIMARY KEY,
  message text NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp
);

CREATE VIRTUAL TABLE fts_files USING fts5(
  body,
  tokenize=unicode61
);

CREATE TABLE files(
  path text PRIMARY KEY,
  title text NOT NULL DEFAULT '',
  tags text NOT NULL DEFAULT ',',
  -- sqlite3def does not support foreign key constraints to virtual tables.
  fts_file_id integer NOT NULL, -- REFERENCES fts_files(rowid),
  modified_at datetime NOT NULL,
  size integer NOT NULL DEFAULT 0,
  updated_at datetime NOT NULL DEFAULT current_timestamp
);

CREATE TABLE users(
  id integer PRIMARY KEY,
  username text NOT NULL,
  password text NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp
);
