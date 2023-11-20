FROM golang:1.20.10-alpine3.18 as builder
WORKDIR /app
ADD ../.. .

RUN CGO_ENABLED=0 GOOS=linux go build -o node-latency-server cmd/server/main.go

FROM scratch as final
COPY --from=builder /app/node-latency-server .
EXPOSE 8080

CMD ["./node-latency-server"]