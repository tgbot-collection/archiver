FROM golang:alpine as builder

RUN apk update && apk add tzdata git make ca-certificates && \
git clone https://github.com/tgbot-collection/archiver /build && \
cd /build && make static


FROM scratch

ENV TZ=Asia/Shanghai
COPY --from=builder /build/archiver /archiver
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /

ENTRYPOINT ["/archiver"]

# docker run -d --restart=always -e TOKEN="FXI" bennythink/archiver