# syntax=docker/dockerfile:1
# Stage 1
FROM golang:1.22.0-alpine3.19 AS dependency_builder
LABEL authors="raoul.mackle@kultivator.co.nz"

WORKDIR /go/src
ENV GO111MODULE=on

RUN apk update
RUN apk add --no-cache bash ca-certificates git make

COPY go.mod .
COPY go.sum .

RUN go mod download

# Stage 2
FROM dependency_builder AS service_builder

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

ARG SERVICE_NAME

USER nonroot:nonroot

WORKDIR /usr/app

COPY common common
COPY database database
COPY services/$SERVICE_NAME/config.yaml services/$SERVICE_NAME/config.yaml
COPY services/$SERVICE_NAME services/$SERVICE_NAME
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o bin services/$SERVICE_NAME/*.go

# Stage 3
FROM alpine:latest

ARG SERVICE_NAME
ARG BUILD_NUMBER

RUN apk update
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root

ENV WORKDIR=services/$SERVICE_NAME
ENV BUILD_NUMBER=$BUILD_NUMBER

COPY --from=service_builder /usr/app/bin bin
COPY --from=service_builder /usr/app/services/$SERVICE_NAME/config.yaml config.yaml
COPY --from=service_builder /usr/app/services/$SERVICE_NAME/templates/ templates/
COPY --from=service_builder /etc/passwd /etc/passwd
COPY --from=service_builder /etc/group /etc/group

EXPOSE 6001

ENTRYPOINT ["/root/bin"]
