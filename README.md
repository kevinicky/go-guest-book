# go-guest-book
![Go](https://github.com/kevinicky/go-guest-book/actions/workflows/go.yml/badge.svg)
![Docker Image](https://github.com/kevinicky/go-guest-book/actions/workflows/docker-image.yml/badge.svg)

Simple Guest Book using Go Postgresql, Redis, Docker, and JWT

- Golang as programming langauge
- Postgresql as main Sql database
- Redis as caching database
- JWT as jwt token authetication
- Docker as container if you dont wanna build the application

## How to run
1. Make sure you have installed below application on your machine:

| Component | Version |
| ----------- | ----------- |
| Golang | 1.19 |
| Postgresql | 14.5 |
| Redis | 6.2.11 |
| Docker | 20.10.23 |

2. Run `git@github.com:kevinicky/go-guest-book.git`
3. Run `PG_DUMP` using 230401_pg_dump.sql
4. Setup `./config/config.yaml` based on your application configuration.
5. If you prefer run the application by build it:
    1.  `go build github.com/kevinicky/go-guest-book/cmd`
    >  if your operating system is windows then it will make `cmd.exe`, otherwise it will make `cmd`
    2.  Run `cmd.exe` or `cmd`
6. If you prefer to use docker image
    1. Run `docker image pull kevinicky/go-guest-book:latest`
    2. Run `docker run -p 8080:8080 kevinicky/go-guest-book`
  
# Quick brief about this application
This is a guest book application. It's use to record guest whenever they visit a place that hosted this appliaction. There is two roles that is:
| Role | Privileges |
| ----------- | ----------- |
| Admin | Create, Update, Delete guest, guest_book, and comment data |
| Non Admin (guest) | Create guest, guest_book, and comment data also Update guest data |

All of this privileges is recorded on `user_matrices` table.
For this application, there is 4 tables:
1. users (record all user data (guest and admin))
2. user_matrices (record all priviledges for admin and guest role)
3. visits (record all visits (or commonly call it s guest_book) by users)
4. threads (record all threads (or commonly call it as comment) by users)

To show flow of this program, please refer to `guest_book.postman_collection`


For admin login:

username : admin@admin.com

password : secret
