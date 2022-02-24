FROM golang:1.17-bullseye as build

WORKDIR /go/src/todo-svc
COPY . /go/src/todo-svc

RUN go get -d -v ./...

RUN go build -o /go/bin/todo-api app/services/todo-api/main.go

FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/todo-api /
CMD ["/todo-api", "--auth-mode=dev"]
