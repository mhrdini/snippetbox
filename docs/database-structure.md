# Database Structure

```sh
snippetbox
|
+-- snippets
     |
     +-- id       INTEGER       NOT NULL PRIMARY KEY AUTO_INCREMENT
     +-- title    VARCHAR(100)  NOT NULL
     +-- content  TEXT          NOT NULL
     +-- created  DATETIME      NOT NULL # has INDEX: idx_snippets_created
     +-- expires  DATETIME      NOT NULL
```
