# Distributed Web Scraper

A scalable, microservice-based web scraper built in Go to extract data from LinkedIn, YouTube, and Instagram. The project leverages modern technologies like Kafka, PostgreSQL, Prometheus, Jaeger, and Consul to ensure reliability, observability, and fault tolerance.

## Table of Contents

* [Project Overview](#project-overview)
* [Architecture](#architecture)
* [Features](#features)
* [Technologies](#technologies)
* [Folder Structure](#folder-structure)
* [Prerequisites](#prerequisites)
* [Setup Instructions](#setup-instructions)
* [Configuration](#configuration)
* [Running the Application](#running-the-application)
* [Monitoring and Observability](#monitoring-and-observability)
* [Testing](#testing)
* [CI/CD](#cicd)
* [Contributing](#contributing)
* [Troubleshooting](#troubleshooting)
* [License](#license)

## Project Overview

The Distributed Web Scraper is designed to scrape data from social media platforms (LinkedIn, YouTube, Instagram) in a scalable and fault-tolerant manner. The application is split into three microservices:

* **Scraper Service**: Scrapes data and produces messages to Kafka.
* **Consumer Service**: Consumes messages and stores them in PostgreSQL.
* **Metrics Service**: Exposes Prometheus metrics.

The system uses Kafka for asynchronous communication, Consul for configuration, and Viper for secrets. Robust error handling and observability tools are included.

## Architecture

### Scraper Service

* Scrapes data using `chromedp`.
* Implements rate limiting, proxy rotation, circuit breakers.
* Publishes data to Kafka topics.
* Uses OAuth for API access.

### Consumer Service

* Consumes Kafka messages.
* Validates data using JSON schema.
* Stores data in PostgreSQL.
* Periodic backups.

### Metrics Service

* Prometheus metrics for performance and errors.
* Integrated with Alertmanager.

### Supporting Services

* Kafka, PostgreSQL, Consul, Prometheus, Grafana, Jaeger, Alertmanager.

## Features

* Microservice Architecture
* Rate Limiting (10 requests/sec)
* Proxy Rotation
* Circuit Breaker (via `gobreaker`)
* OAuth Authentication
* JSON Schema Validation
* Kafka-based Asynchronous Processing
* Centralized Configuration with Consul
* Secrets Management with Viper
* Observability with Prometheus, Grafana, Jaeger, Alertmanager
* Periodic DB Backups
* Correlation ID Logging with Zap
* CI/CD with GitHub Actions
* Unit Testing

## Technologies

* Go (v1.21)
* Kafka (Confluent)
* PostgreSQL (v15)
* Consul
* Viper
* Prometheus
* Grafana
* Jaeger
* Alertmanager
* chromedp
* gobreaker
* OpenTelemetry
* Zap Logger
* Docker & Docker Compose

## Folder Structure

```
distributed-web-scraper/
├── services/
│   ├── scraper/
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── scraper/
│   │   │   │   ├── linkedin/
│   │   │   │   ├── instagram/
│   │   │   │   ├── youtube/
│   │   │   ├── kafka/
│   │   │   ├── auth/
│   │   │   ├── logger/
│   │   │   ├── tracing/
│   │   ├── Dockerfile
│   │   ├── .env
│   ├── consumer/
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── kafka/
│   │   │   ├── storage/
│   │   │   ├── validation/
│   │   │   ├── logger/
│   │   │   ├── tracing/
│   │   ├── Dockerfile
│   │   ├── .env
│   ├── metrics/
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── metrics/
│   │   │   ├── logger/
│   │   ├── Dockerfile
│   │   ├── .env
├── sql/
├── docker/
│   ├── prometheus.yml
│   ├── alerts.yml
│   ├── consul_config.json
│   ├── docker-compose.yml
├── .github/workflows/
├── go.mod
├── go.sum
├── README.md
```

## Prerequisites

* Docker & Docker Compose
* Go 1.21 (optional)
* Git
* OAuth credentials for LinkedIn, YouTube, Instagram
* Proxy servers (recommended)

## Setup Instructions

### Clone the Repo

```sh
git clone https://github.com/YourOrg/distributed-web-scraper.git
cd distributed-web-scraper
```

### Configure Secrets

Edit `services/scraper/.env`:

```
LINKEDIN_CLIENT_ID=your-linkedin-client-id
LINKEDIN_CLIENT_SECRET=your-linkedin-client-secret
YOUTUBE_CLIENT_ID=your-youtube-client-id
YOUTUBE_CLIENT_SECRET=your-youtube-client-secret
INSTAGRAM_CLIENT_ID=your-instagram-client-id
INSTAGRAM_CLIENT_SECRET=your-instagram-client-secret
```

Edit `services/consumer/.env`:

```
POSTGRES_URL=postgres://user:password@postgres:5432/scraper?sslmode=disable
```

Leave `services/metrics/.env` empty.

### Configure Consul

Edit `docker/consul_config.json`:

```json
{
  "scraper/config": {
    "KAFKA_BROKERS": ["kafka:9092"],
    "SCRAPE_INTERVAL": 300,
    "PROXY_LIST": ["proxy1:port", "proxy2:port"]
  }
}
```

Replace with your proxy list.

### Install Dependencies

```sh
go mod tidy
```

## Configuration

* **Scraper**: `.env`, Consul (rate limit, circuit breaker, scrape interval)
* **Consumer**: `.env`, Consul (PostgreSQL, Kafka)
* **Metrics**: Consul (no secrets)
* **Kafka Topics**: `linkedin_data`, `instagram_data`, `youtube_data`
* **PostgreSQL Schema**: `scraped_data` (JSONB column)

## Running the Application

### Start with Docker Compose

```sh
docker-compose -f docker/docker-compose.yml up -d
```

### Verify

* Scraper: `docker logs <scraper-container>`
* Consumer: `docker logs <consumer-container>`
* Metrics: [http://localhost:9090/metrics](http://localhost:9090/metrics)

### Stop Services

```sh
docker-compose -f docker/docker-compose.yml down
```

## Monitoring and Observability

* **Prometheus**: [http://localhost:9090](http://localhost:9090)

  * Metrics: `scraper_requests_total`, `scraper_duration_seconds`, `scraper_errors_total`
* **Grafana**: [http://localhost:3000](http://localhost:3000) (admin/admin)
* **Jaeger**: [http://localhost:16686](http://localhost:16686)
* **Alertmanager**: [http://localhost:9093](http://localhost:9093)
* **Consul**: [http://localhost:8500](http://localhost:8500)



## CI/CD

GitHub Actions pipeline:

* Build Go application
* Run tests
* Build Docker images

Extend with deployment steps as needed.

## Contributing

1. Fork the repo
2. Create a branch (`git checkout -b feature/your-feature`)
3. Commit (`git commit -m "Add your feature"`)
4. Push (`git push origin feature/your-feature`)
5. Open a pull request

Follow Go standards and write tests.

## Troubleshooting

* **Constructor Errors**:

  * Match imports with `go.mod`
  * Run `go mod tidy`
  * Check file paths for scrapers

* **Kafka Issues**:

  * Verify `KAFKA_BROKERS` in Consul
  * Ensure Kafka/Zookeeper are running

* **DB Errors**:

  * Check `POSTGRES_URL`
  * Inspect PostgreSQL logs

* **Scraping Failures**:

  * Update selectors and URLs
  * Validate proxy servers

* **Metrics Missing**:

  * Ensure Metrics service is live
  * Validate Prometheus config

## License

MIT License. See `LICENSE` file for details.
