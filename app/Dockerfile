FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o rolevod main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/rolevod /app/rolevod

ENTRYPOINT ["/app/rolevod"]
