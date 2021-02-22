FROM golang:latest as builder

WORKDIR sciienceobjectsdb

COPY . .

CMD CGO_ENABLED=0 GOOS=linux go test ./...