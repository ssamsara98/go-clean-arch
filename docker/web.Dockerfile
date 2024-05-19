FROM golang:1.20-alpine

# Required because go requires gcc to build
RUN apk add build-base git inotify-tools
RUN echo $GOPATH
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY . .

CMD sh /app/docker/run.sh
