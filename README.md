# Snippetbox

A tutorial project written in Go, completed by following the Let's Go book.

## Table of Contents:

- [Project Structure](./docs/project-structure.md)

## Development Mode

**Requirements:**

- Go v1.20
- [air](https://github.com/cosmtrek/air) for hot reloading
- [MySQL DB Setup](https://github.com/mhrdini/snippetbox/blob/main/docs/database.md)
- [TLS Certificate Generation](https://github.com/snippetbox/blob/main/docs/tls.md)

**Steps:**

```bash
go mod tidy
cd cmd/web
air
```

To kill a process at some port:

```bash
kill $(lsof -t -i:PORT)
```
