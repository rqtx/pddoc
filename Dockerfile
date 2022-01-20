# builder image
FROM golang:1.17.6-alpine AS builder

RUN apk add --update --no-cache make ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /build/pddoc
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make go-build

# generate clean, final image for end users
FROM amazon/aws-cli:2.4.12
COPY --from=builder /build/pddoc/bin/pddoc /usr/local/bin/
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "pddoc" ]