FROM golang:1.21

WORKDIR /app
COPY services/metrics .
RUN go mod download
RUN go build -o metrics ./cmd

CMD ["./metrics"]