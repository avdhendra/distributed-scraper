FROM golang:1.21

WORKDIR /app
COPY services/consumer .
RUN go mod download
RUN go build -o consumer ./cmd

CMD ["./consumer"]