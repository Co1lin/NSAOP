FROM golang:latest AS build_base

ENV GOPROXY=https://goproxy.cn,direct

# Disable compiling with CGO
# to make sure that the executable file can run on alpine
ENV CGO_ENABLED=0

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Generate swagger docs
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init

# Build the Go app
RUN go build -o ./nsaop-backend .

# --------
# Start fresh from a smaller image
FROM alpine:latest

COPY --from=build_base /app/nsaop-backend /app/nsaop-backend
COPY --from=build_base /app/config /app/config
COPY --from=build_base /app/docs /app/docs

WORKDIR /app

# This container exposes port 10080 to the outside world
EXPOSE 10080

# Run the binary program produced by `go install`
CMD ["./nsaop-backend"]
