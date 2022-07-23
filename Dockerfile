FROM golang:1.17.4-alpine as builder

RUN apk update && apk add --no-cache make ca-certificates tzdata && mkdir /build
COPY go.mod /build
RUN cd /build && go mod download
COPY . /build
RUN cd /build && make static


FROM alpine
RUN apk update && apk add --no-cache chromium-chromedriver

ENV TZ=Asia/Shanghai
COPY --from=builder /build/archiver /archiver
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /

CMD ["/archiver"]

# docker run -d --restart=always -e TOKEN="FXI" bennythink/archiver