FROM golang:alpine as builder

RUN apk update && apk add git make ca-certificates && \
git clone https://github.com/tgbot-collection/archiver /build && \
cd /build && make static


FROM scratch

COPY --from=builder /build/archiver /archiver
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /

ENTRYPOINT ["/archiver"]

# docker run -d --restart=always -e TOKEN="FXI" bennythink/archiver