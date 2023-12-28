# Build stage
FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /restApi

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o GolandRestApi ./cmd/server

# Final stage - all the code
EXPOSE 8080

CMD ["./GolandRestApi"]

# Final stage - only with the executable
#FROM alpine:3.19.0
#
#WORKDIR /root/
#
#COPY --from=builder /app/GolandRestApi .
#
#EXPOSE 8080
#
#CMD ["./GolandRestApi"]
