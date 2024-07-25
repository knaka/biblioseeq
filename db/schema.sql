CREATE TABLE logs (
  id integer PRIMARY KEY,
  message text NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp
);

CREATE VIRTUAL TABLE documents USING fts5(
  title,
  body,
  tokenize=unicode61
);

CREATE TABLE files(
  path text PRIMARY KEY,
  document_id int NOT NULL REFERENCES documents (docid),
  modified_at datetime NOT NULL,
  size integer NOT NULL DEFAULT 0,
  updated_at datetime NOT NULL DEFAULT current_timestamp
);

CREATE INDEX index_139b076 ON files(document_id);
