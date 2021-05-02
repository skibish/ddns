# To build: docker build -t ddns .
# To run:   docker run -v /path/to/config.yml:/config/.ddns.yml ddns -conf-file /config/.ddns.yml
# Or if your .ddns.yml is in the current working directory and is named .ddns.yml
# docker run -v ${PWD}:/config ddns
FROM golang:1.16.3-alpine as builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy everything in and do the go build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddns .

# Now create a new stage and only copy the binary we need (keeps the container small)
FROM scratch
COPY --from=builder /app/ddns /

ENTRYPOINT ["/ddns"]
