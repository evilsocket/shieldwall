FROM golang:alpine as builder

RUN apk update && apk add --no-cache make

# download, cache and install deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# copy and compiled the app
COPY . .
RUN make clean
RUN make api

# start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/_build/shieldwall-api .

EXPOSE 8666
ENTRYPOINT ["./shieldwall-api"]