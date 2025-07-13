# Backend Rewrite: morigin-server in Go

This document outlines the architecture and design decisions for the rewrite of the `morigin-server` from Java to Go.

## 1. Project Overview

The primary goal of this project is to rewrite the existing `morigin-server` from Java to Go. The new implementation will use a PostgreSQL database instead of the current database.

## 2. Architecture

We will be using a Hexagonal Architecture (also known as Ports and Adapters). This architecture promotes a clear separation of concerns, making the application more modular, testable, and maintainable. The architecture is divided into three main layers:

*   **Domain:** This is the core of the application and contains the business logic and models. It is completely independent of any external frameworks or libraries.
*   **Application:** This layer orchestrates the calls to the domain layer and is responsible for the application's use cases.
*   **Infrastructure:** This layer contains all the external dependencies, such as the database, web framework, and other third-party libraries.

## 3. Directory Structure

The directory structure will be organized as follows:

```
backend/
├── cmd/
│   └── server/
│       └── main.go         # Entry point of the application
├── configs/                  # Configuration files
│   └── config.yaml
├── migrations/               # Database migrations
├── internal/
│   ├── app/                  # Application layer
│   ├── domain/               # Domain layer
│   │   └── user/             # User model
│   └── infrastructure/       # Infrastructure layer
│       ├── config/           # Configuration loading
│       ├── persistence/      # Database connection and repositories
│       └── web/              # Web framework and handlers
│           └── middleware/   # Middleware for the web framework
├── pkg/
│   └── auth/                 # Authentication (JWT)
├── go.mod
└── go.sum
```

## 4. Database

We will be using PostgreSQL as the database. We will use [GORM](https://gorm.io/) as the ORM to interact with the database.

## 5. API

The API will be a RESTful API. We will use the [Gin](https://gin-gonic.com/) web framework to build the API.

## 6. Authentication

Authentication will be handled using JSON Web Tokens (JWT).

## 7. Configuration

Application configuration will be managed using [Viper](https://github.com/spf13/viper). The configuration will be stored in a `config.yaml` file.

## 8. Testing

We will be writing unit tests for the domain and application layers. We will also write integration tests for the API endpoints.

## 9. Dependencies

*   **Web Framework:** Gin
*   **ORM:** GORM
*   **Database Driver:** pgx
*   **Configuration:** Viper
*   **Authentication:** golang-jwt
