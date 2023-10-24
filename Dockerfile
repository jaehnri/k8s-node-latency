FROM golang:1.20.10-alpine3.18 as builder

WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-node-latency

FROM scratch as final

COPY --from=builder /app/k8s-node-latency .

EXPOSE 8080

CMD ["./k8s-node-latency"]