# Snippetbox

A tutorial project written in Go, completed by following the Let's Go book.

## Table of Contents:

- [Project Structure](./docs/project-structure.md)

## Development Mode

**Requirements:**

- Go v1.20
- [air](https://github.com/cosmtrek/air) for hot reloading

**Database Setup:**

```bash
$ sudo mysql
mysql>

# Create snippetbox DB
mysql> CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
mysql> USE snippetbox;

# Create snippets table with index on create date
mysql> CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
);
mysql> CREATE INDEX idx_snippets_created ON snippets(created);

# Create 'web' user with 'web' password
mysql> CREATE USER 'web'@'localhost';
mysql> GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
mysql> ALTER USER 'web'@'localhost' IDENTIFIED BY 'web';

# Insert dummy records
mysql> INSERT INTO snippets (title, content, created, expires) VALUES (
  'An old silent pond',
  'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n- Matsuo Bashō', UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
mysql> INSERT INTO snippets (title, content, created, expires) VALUES (
  'Over the wintry forest',
  'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
mysql> INSERT INTO snippets (title, content, created, expires) VALUES (
  'First autumn morning',
  'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n- Murakami Kijo',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
```

**Steps:**

```bash
go mod tidy
cd cmd/web
air
```
