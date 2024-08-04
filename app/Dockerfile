FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o rolevod app/main.go

FROM alpine

WORKDIR /app

COPY app/reference.yaml .

COPY --from=builder /build/rolevod .

ENTRYPOINT ["/app/rolevod"]
