FROM golang:1.20.6 AS builder
WORKDIR /go/src/export

RUN apt-get update
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:3.16
WORKDIR /app/

COPY --from=builder /go/src/export/bin/export .
COPY --from=builder /go/src/export/cmd/.env .
COPY --from=builder /go/src/export/database/migrations ./database/migrations

EXPOSE 80
EXPOSE 443
CMD ["./export"]