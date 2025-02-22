FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /go/bin/app ./server.go

FROM gcr.io/distroless/base-debian12

COPY --from=builder /go/bin/app /app

EXPOSE 8080

CMD ["/app"]