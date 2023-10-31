FROM golang:1.20.10-alpine3.18 as builder
WORKDIR /app
ADD ../.. .

RUN CGO_ENABLED=0 GOOS=linux go build -o node-latency-client cmd/client/main.go

FROM scratch as final
COPY --from=builder /app/node-latency-client .

CMD ["./node-latency-client"]