# CYBER GAME

Gin is a web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. It also provides a robust set of features for building web applications and APIs.

## Prerequisites

Before you can start using Gin, you need to have the following installed:

- Go version 1.20 or higher
- .env file
- mgration.env file for test migration database

## Getting Started

```
go run .
```

## How to Migration MySQL

Create a file name migrate.env to root path

```
DB_USER = user
DB_PASS = pass
DB_HOST = host
DB_PORT = 3360
DB_NAME = cybergame
```

Run migrate
```
go run migration/migrate.go up
go run migration/migrate.go down 1
```

## Example APIs

| METHOD | URL | TOKEN |
|--------|-----|-------|
| GET | http://localhost:3000/login | asdasd |
| POST | http://localhost:3000/login | asdasd |
| POST | http://localhost:3000/register | REQUIRED |
| PATCH | http://localhost:3000 | REQUIRED |
| DELETE | http://localhost:3000 | REQUIRED |

## Document

URL: http://localhost:3000/swagger/index.html