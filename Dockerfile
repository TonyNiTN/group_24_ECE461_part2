FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/main /app/main

EXPOSE 8080

ENV PORT 8080

ENV LOG_LEVEL 2
ENV LOG_FILE mylog.log

CMD ["/app/main"]