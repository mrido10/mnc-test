## Soal 1
Before run the apps `test1`, you can change directory in `~/test1`, then run command 
```shell
go mod tidy

go run main.go
```

## Soal 2
Before run the apps `test2`, you can change directory in `~/test2`, then run command
```shell
go mod tidy

# run container postgresql and redis
docker compose up -d

# set environtment variable
export PG_DB=mnc-db
export PG_HOST=localhost
export PG_PASS=thepassword
export PG_PORT=5433
export PG_USER=idopsql
export REDIS_HOST=localhost
export REDIS_PORT=6379

# run apps
go run main.go

```

