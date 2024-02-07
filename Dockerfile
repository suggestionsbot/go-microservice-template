FROM golang:1.22-alpine AS build

WORKDIR /go/src/microservice

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o microservice .

FROM alpine:latest

WORKDIR /opt/microservice

COPY --from=build /go/src/microservice .

EXPOSE 3000

CMD ["./main"]