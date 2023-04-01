# go-guest-book
Simple Guest Book using Go Postgresql, Redis, and JWT

- Golang as programming langauge
- Postgresql as main Sql database
- Redis as caching database
- JWT as jwt token authetication

## How to run
1. Make sure you have installed below application on your machine:

| Component | Version |
| ----------- | ----------- |
| Golang | 1.19 |
| Postgresql | 14.5 |
| Redis | 6.2.11 |

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
  
    
