FROM golang:alpine as build

RUN apk update
RUN apk upgrade

# RUN apk add --update go=1.8.3-r0 gcc=6.3.0-r4 g++=6.3.0-r4
RUN apk add --no-cache git gcc g++ make

WORKDIR /

# Cache dependencies in seperate layer
COPY go.mod go.sum makefile ./
RUN make install

COPY . .

RUN make build

# Build the default development membrane server
FROM alpine
# FIXME: Build these in a build stage during the docker build
# for now will just be copied post local build
# and execute these stages through a local shell script
COPY --from=build ./bin/api-gateway /api-gateway
RUN chmod +rx /api-gateway

ENTRYPOINT [ "/api-gateway" ]