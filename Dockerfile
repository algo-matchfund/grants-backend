FROM quay.io/goswagger/swagger AS swagger
WORKDIR /go/src/github.com/algo-matchfund/grants-backend
COPY . .

# prepare api files
RUN /usr/bin/swagger generate server -A grants-program -t ./gen --exclude-main  -f ./open-api-specifications/grants-program/grants-program.yaml -P models.Principal --skip-tag-packages

FROM alpine:latest as builder

RUN apk add --no-cache \
        go \
        gcc \
        musl-dev

# copy generated api files from swagger container
WORKDIR /go/src/github.com/algo-matchfund/grants-backend
COPY --from=swagger /go/src/github.com/algo-matchfund/grants-backend .

# build grants-backend
RUN go get -d -v ./cmd/grants-backend
RUN go build -v -a -o ./build/grants-backend ./cmd/grants-backend

FROM alpine:latest

RUN apk add --no-cache \
        libstdc++      \
        libgcc

# copy the binary
WORKDIR /home/
COPY --from=builder /go/src/github.com/algo-matchfund/grants-backend/build/grants-backend .
RUN mkdir -p db/migration
COPY ./db/migration/* db/migration/
COPY ./ci/certs certs

# expose and run server
EXPOSE 8080/tcp
ENTRYPOINT ["./grants-backend", "--config", "config.yml"]

