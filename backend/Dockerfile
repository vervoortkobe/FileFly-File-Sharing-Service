FROM golang:1.22 AS builder
WORKDIR /backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:3.20
WORKDIR /backend
COPY --from=builder /backend/backend .
EXPOSE 8080
CMD ["./backend"]
