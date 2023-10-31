FROM golang:1.20.10-alpine3.18 as builder
WORKDIR /app
ADD ../.. .

RUN CGO_ENABLED=0 GOOS=linux go build -o client-node-latency cmd/client/main.go

FROM scratch as final
COPY --from=builder /app/client-node-latency .

CMD ["./client-node-latency"]