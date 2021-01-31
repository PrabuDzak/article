# Article

A weekend Go exercise project

# App

## `restapi`

REST API application. Accept and respond in JSON

- `GET /articles`
  - query paremeter: `author`, `keyword`, `limit`, `offset`
- `POST /articles`
  - body parameter: 
    ```json
      {
        "author": "string,required",
        "title": "string,required",
        "body": "string,required"
      }
    ```


# Require

- go 1.13
- docker
- docker-compose

# How to

## Test

```sh
make test
```

## Run

```sh
docker-compose up -d  # prepare env. it takes time
cp env.sample .env    # create env var file
make migrate          # load/migrate database schema
make mapping          # apply index mappings
make compile          # compile 
make run              # run
```

## Run Acceptence Test

```sh
make acceptence
```

## Run in Docker

```
make build
docker run --network host --env-file .env prabudzak/article:latest
```
