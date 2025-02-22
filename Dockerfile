FROM alpine:latest AS builder

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ddns /

ENTRYPOINT ["/ddns"]
