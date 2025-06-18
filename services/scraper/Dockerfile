FROM golang:1.21

WORKDIR /app
COPY services/scraper .
RUN go mod download
RUN go build -o scraper ./cmd

CMD ["./scraper"]