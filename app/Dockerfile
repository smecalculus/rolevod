FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o rolevod app/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/rolevod .

COPY ./app/reference.yaml .

ENTRYPOINT ["/app/rolevod"]
