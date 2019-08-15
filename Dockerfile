# To build: docker build -t ddns .
# To run:   docker run -v /path/to/config.yml:/config/.ddns.yml ddns -conf-file /config/.ddns.yml
# Or if your .ddns.yml is in the current working directory and is named .ddns.yml
# docker run -v ${PWD}:/config ddns
FROM golang:1.12.8-alpine as builder

RUN apk update && \
  apk upgrade && \
  apk add git

ENV GO111MODULE=on
WORKDIR /app

# Cache dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything in and do the go build
COPY . .
RUN go build

# Now create a new stage and only copy the binary we need (keeps the container small)
FROM golang:1.12.8-alpine
COPY --from=builder /app/ddns /app/

ENTRYPOINT ["/app/ddns"]
