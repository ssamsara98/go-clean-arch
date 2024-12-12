ARG IMAGE=golang:1.21-alpine

FROM ${IMAGE} as base

ENV PORT=8080
ENV ENVIRONMENT=production
ENV LOG_LEVEL=info
ENV LOG_OUTPUT=server.log
ENV DB_TYPE=postgres
ENV MAX_MULTIPART_MEMORY=10485760
ENV JWT_ACCESS_SECRET=9f4b7b1e0c2c166aa1733cdcf1f1c1a3f2ded287b5d28c2194378273f5530bd8
ENV JWT_REFRESH_SECRET=9bc0195061abea1f9461ba84412c9d8819594a1b166b0e1996ca097414224ced
ENV ACCESS_TOKEN_DURATION=24h
ENV REFRESH_TOKEN_DURATION=720h

EXPOSE 8080
EXPOSE 4000

# Required because go requires gcc to build
RUN apk update && apk upgrade && apk add --no-cache build-base bash git inotify-tools make
RUN go install github.com/rubenv/sql-migrate/...@latest

WORKDIR /app
COPY --link go.mod go.sum ./
RUN go mod download

FROM base as build
RUN mkdir -p /tmp/build
COPY --from=base --link /app/go.mod /app/go.sum /tmp/build/
RUN cd /tmp/build && go mod tidy && go mod verify
COPY --link . .
RUN go build -buildvcs=false -o main

FROM base AS release
COPY --from=build /app/migration migration
COPY --from=build /app/seeders seeders
COPY --from=build /app/go.mod /app/go.sum /app/make /app/Makefile /app/dbconfig.yaml /app/template.production.env /app/main ./
COPY --link ./home/.profile /root/

RUN mv template.production.env .env
CMD ["./main", "app:serve"]
