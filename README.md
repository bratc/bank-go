This repository holds code and test result that were made for the article https://dzone.com/articles/java-vs-go-microservices-load-testing-rematch

This code was copied from the https://github.com/nikitsenka/bank-go which is a part of original article https://dzone.com/articles/java-vs-go-multiple-users-load-test-1

======= From original README.md =======
Run tests locally
```
  docker pull postgres
  docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=test1234 -d postgres
  go test ./bank
```
Build docker image
```
  docker build --no-cache -t bank-go .

```
Run app in Docker
with external postgres
```
  docker run --name bank-go -p 8000:8000 -e POSTGRES_HOST=${host} -d bank-go
```
  or create both postgres and bank-go containers and run
```
  docker-compose up -d --force-recreate
```
