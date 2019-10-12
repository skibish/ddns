# To build: docker build -t ddns .
# To run:   docker run -v /path/to/config.yml:/config/.ddns.yml ddns -conf-file /config/.ddns.yml
# Or if your .ddns.yml is in the current working directory and is named .ddns.yml
# docker run -v ${PWD}:/config ddns
FROM golang:1.13.1-alpine as builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy everything in and do the go build
COPY . .
RUN go build -v

# Now create a new stage and only copy the binary we need (keeps the container small)
FROM alpine:3.10.2
COPY --from=builder /app/ddns /

ENTRYPOINT ["/ddns"]
