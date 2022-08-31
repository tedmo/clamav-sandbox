FROM golang:1.18 AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/app ./cmd/app

FROM alpine:3.16 AS final

EXPOSE 8080

RUN apk update \
    && apk add --no-cache pv ca-certificates clamav clamav-libunrar \
    && apk add --upgrade apk-tools libcurl openssl busybox \
    && rm -rf /var/cache/apk/*

# Copy app
COPY --from=build /usr/local/bin/app /usr/local/bin/app

# Temp directory to hold files to scan
RUN mkdir -p /temp/scan

ENTRYPOINT ["app"]