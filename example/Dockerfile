FROM registry.trendyol.com/platform/base/image/golang:1.21.3-alpine as build
ENV TZ=Europe/Istanbul
WORKDIR /app

RUN go get -u github.com/esdonmez/deneme-main

WORKDIR /app

RUN profile_downloader

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -pgo=merged.pprof -ldflags="-w -s" -o /app/buybox

# Switch to a small base image
FROM registry.trendyol.com/platform/base/image/appsec/chainguard/static/library:lib-20230201

# copy app to bin directory, and set it as entrypoint
WORKDIR /app
# Get the TLS CA certificates from the build container, they're not provided by busybox.
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Istanbul
ENV ZONEINFO=/zoneinfo.zip
COPY --from=build /app/buybox /app/buybox

EXPOSE 8082

ENTRYPOINT ["/app/buybox"]
