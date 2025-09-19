FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /sca ./cmd/main.go

FROM alpine:3.18.0
COPY --from=builder /sca /sca
CMD [ "/sca" ]
