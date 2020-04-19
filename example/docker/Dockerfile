FROM golang:1.14.2-alpine3.11 as build-env
# All these steps will be cached
RUN apk add --update --no-cache \
        git ca-certificates tzdata && update-ca-certificates
RUN mkdir /app
WORKDIR /app
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app cmd/main.go

# <- Second step to build minimal image
FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=build-env /go/bin/app /go/bin/app

# barcode font
COPY ./luxisr.ttf /go/bin/

#ENV
#Default to Asia/Seoul
ENV TZ Asia/Seoul

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/app"]
