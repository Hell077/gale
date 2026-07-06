FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY sdk ./sdk

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gale ./cmd

FROM alpine:3.22

RUN addgroup -S gale && adduser -S gale -G gale

USER gale
WORKDIR /app

COPY --from=build /out/gale /usr/local/bin/gale

ENV GALE_HOST=0.0.0.0
ENV GALE_PORT=9000

EXPOSE 9000

ENTRYPOINT ["gale"]
