# Database

## Structure

```sh
snippetbox
|
+-- snippets
|     |
|     +-- id        INTEGER        NOT NULL PRIMARY KEY AUTO_INCREMENT
|     +-- title     VARCHAR(100)   NOT NULL
|     +-- content   TEXT           NOT NULL
|     +-- created   DATETIME       NOT NULL # has INDEX: idx_snippets_created
|     +-- expires   DATETIME       NOT NULL
|
+-- sessions
|     |
|     +-- token     CHAR(43)       PRIMARY KEY
|     +-- data      BLOB           NOT NULL
|     +-- expiry    TIMESTAMP(6)   NOT NULL
|
+-- users
      |
      +-- id                INTEGER       NOT NULL PRIMARY KEY AUTO_INCREMENT
      +-- name              VARCHAR(255)  NOT NULL
      +-- email             VARCHAR(255)  NOT NULL
      +-- hashed_password   CHAR(60)      NOT NULL
      +-- created           DATETIME      NOT NULL

test_snippetbox
|
+-- snippets

```

## Setup

After installing MySQL on your machine:

```bash
$ sudo mysql # or mysql -u root -p
# when prompted for a password, enter your password and/or just press Enter
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

# Create sessions table
mysql> CREATE TABLE sessions (
  token CHAR(43) PRIMARY KEY,
  data BLOB NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
);

# Create users with unique constraint on email column
mysql> CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created DATETIME NOT NULL
);

mysql> ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

### Testing Database

```bash
# Create test database and user
mysql> CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

mysql> CREATE USER 'test_web'@'localhost';
mysql> GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_snippetbox.* TO 'test_web'@'localhost';
mysql> ALTER USER 'test_web'@'localhost' IDENTIFIED BY 'test_web';
```

Also used:

- [setup.sql](https://github.com/mhrdini/snippetbox/blob/main/internal/models/testdata/setup.md)
- [teardown.sql](https://github.com/mhrdini/snippetbox/blob/main/internal/models/testdata/teardown.md)
