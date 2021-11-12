FROM golang:1.15-alpine AS builder
WORKDIR /chat
ENV GO111MODULE=on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o chat ./cmd/web/*.go

FROM alpine:3.12
WORKDIR /chat
COPY --from=builder /chat/chat ./
COPY --from=builder  /chat/html/ ./html/

RUN chmod +x /chat

ENTRYPOINT ["/chat/chat"]