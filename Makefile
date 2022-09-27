#!make

-include .env

mock:
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4
	go generate ./...

test:
	go test -cover `go list ./... | grep -Ev 'main|mock|testing'`	

compile:
	CGO_ENABLED=0 GOOS=linux go build -o ./_output/restapi ./app/restapi/main/main.go
	CGO_ENABLED=0 GOOS=linux go build -o ./_output/testing ./app/testing/main/main.go

build:
	docker build --no-cache -t prabudzak/article:latest -f Dockerfile .

build-docker:
	docker build --no-cache -t prabudzak/article:latest -f Dockerfile.build .

run:
	CGO_ENABLED=0 GOOS=linux go build -o ./_output/restapi ./app/restapi/main/main.go
	CGO_ENABLED=0 GOOS=linux go build -o ./_output/testing ./app/testing/main/main.go
	./_output/restapi

acceptence:
	./_output/testing

migrate:
	@ ./bin/migrate -path="./db/migration" -database="mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" up

mapping:
	curl -i -X PUT $(ELASTICSEARCH_URL)/$(ELASTICSEARCH_ARTICLE_INDEX) -H "Content-Type: application/json" --data-binary "@db/indexmapping/article.json"
