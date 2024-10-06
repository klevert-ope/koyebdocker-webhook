# Koyeb Docker Webhook Handler

This is a webhook handler for Koyeb that listens for Docker Hub webhook events and triggers redeployment of services.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Logging](#logging)
- [Health Check](#health-check)
- [Contributing](#contributing)
- [License](#license)

## Overview

This application is designed to handle incoming webhooks from Docker Hub and trigger redeployments on Koyeb for specified services. It listens on a specified port, processes the webhook payload, and makes API calls to Koyeb to redeploy the services based on the received image tags.

## Features

- Handles Docker Hub webhooks
- Triggers redeployment on Koyeb
- Configurable via environment variables
- Graceful shutdown support
- Health check endpoint

## Getting Started

### Prerequisites

- Go 1.23+
- Docker
- Koyeb API token
- Docker Hub webhook configured

### Installation

Clone the repository:

```sh
git clone https://github.com/klevert-ope/koyebdocker-webhook.git
cd koyeb-docker-webhook-handler
```

Build the Docker image:

```sh
docker build -t koyeb-webhook-handler .
```

### Configuration

The application is configured via environment variables. You need to set the following environment variables:

- `KOYEB_API_TOKEN`: Your Koyeb API token
- `SERVICE_1_ID`, `SERVICE_1_IMAGE`: The Koyeb service ID and corresponding Docker image for the first service
- `SERVICE_2_ID`, `SERVICE_2_IMAGE`: The Koyeb service ID and corresponding Docker image for the second service
- ... and so on for each service you want to configure

### Running the Application

Run the Docker container with the necessary environment variables:

```sh
docker run -d -p 8080:8080 --name koyeb-webhook-handler \
  -e KOYEB_API_TOKEN=your_koyeb_api_token \
  -e SERVICE_1_ID=your_service_1_id \
  -e SERVICE_1_IMAGE=your_service_1_image \
  -e SERVICE_2_ID=your_service_2_id \
  -e SERVICE_2_IMAGE=your_service_2_image \
  koyeb-webhook-handler
```

## Usage

Configure your Docker Hub repository to send webhooks to your application's URL (e.g., `http://your-server-ip:8080/webhook`). When a new image is pushed, Docker Hub will send a webhook to your application, which will then trigger a redeployment on Koyeb for the configured services.

## Endpoints

- `POST /webhook`: Endpoint to handle Docker Hub webhooks
- `GET /health`: Health check endpoint

## Logging

The application logs important events such as incoming webhook payloads, service redeployment triggers, and errors. Logs are output to the console and can be viewed using Docker logs:

```sh
docker logs koyeb-webhook-handler
```

## Health Check

The application provides a health check endpoint at `/health`. It can be used to check if the application is running and ready to handle requests.

Example health check request:

```sh
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## Contributing

Contributions are welcome! Please fork the repository and open a pull request with your changes. Ensure that your code follows the best practices and includes tests where applicable.

## License

This project is licensed under the MIT License.
```

Replace `https://github.com/klevert-ope/koyebdocker-webhook.git` with your actual GitHub repository URL and adjust any other details specific to your project. This README provides a comprehensive guide on what the project is, how to set it up, and how to use it.