
# Stratflow

## About

Stratflow is a trading strategy builder for designing, configuring, and testing strategy hypotheses.

The core idea is to keep strategy configuration outside of the codebase. Each strategy can have multiple parameters, and each parameter can have multiple values. This makes it possible to build relationship graphs inside a strategy, test different parameter combinations, and add new hypotheses without hardcoding them directly into the strategy implementation.

Stratflow is intended to be used as a foundation for strategy research, configuration management, and future automation around trading strategy construction.


## Architecture

Stratflow is built as a microservice-based application. The current architecture includes a gRPC-based Identity service responsible for authentication and user management, a Strategy Registry service responsible for storing and managing strategies, parameters, and parameter values, a React-based web interface for application management, PostgreSQL as the main database, and nginx as an API gateway and reverse proxy.



## How to start

Clone the repository:

```bash
git clone https://github.com/stratflow-labs/stratflow.git
cd stratflow
```

Create a local environment file. Change values if needed:

```bash
cp .env.example .env
```

Start the application:

```bash
docker compose up -d
```

Run migrations:

```bash
make migrations
```
Create admin:

```bash
make admin
```

## Branches & Contributing

- The main branch is the default branch used for production deployments. Changes to this branch are made from the staging branch once a version is ready for community use.

- The staging branch is used for pre-release testing. It is stable enough for testing but not yet ready for production deployment.

- The develop branch is used for development and is the default branch for contributions.
