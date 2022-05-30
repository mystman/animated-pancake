FROM golang:latest as builder
RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build_release


FROM alpine:latest AS prod
COPY --from=builder /app/bin/service .
CMD ["./service"]