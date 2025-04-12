# Modular Go Application

A modular Go application with a web module using Gin and an items module that implements a REST API.

## Project Structure

```
modularApp/
├── cmd/
│   └── api/              # Application entry points
│       └── main.go
├── internal/
│   └── server/           # Internal server implementation
│       └── server.go
├── pkg/
│   ├── web/              # Web module (core HTTP functionality)
│   │   └── module.go     # Module definition and base routing
│   └── items/            # Items module (items REST API)
│       ├── module.go     # Module definition and route handlers
│       ├── models.go     # Data models
│       └── repository.go # Data repository
├── go.mod                # Go module file
└── README.md             # This file
```

## Modular Architecture

This application follows a modular architecture pattern:

- **Web Module**: Core HTTP server and routing functionality using Gin
- **Items Module**: Business logic for item management, depends on Web Module
- **Server**: Orchestrates and wires up all modules
- **Main**: Application entry point that starts the server

Each module is independent and has a clear responsibility, making the codebase more maintainable and extensible.

## Prerequisites

- Go 1.21 or later

## Getting Started

1. Install dependencies:
   ```
   go mod download
   ```

2. Run the application:
   ```
   go run cmd/api/main.go
   ```

3. Test the endpoints:
   ```
   curl http://localhost:8080/api/health
   curl http://localhost:8080/api/v1/ping
   ```

## REST API Endpoints

The application provides a RESTful API for managing items:

### Items Endpoints

- `GET /api/v1/items` - Get all items
- `GET /api/v1/items/:id` - Get item by ID
- `POST /api/v1/items` - Create a new item
  ```json
  {
    "name": "Item name",
    "content": "Item content"
  }
  ```
- `PUT /api/v1/items/:id` - Update an existing item
  ```json
  {
    "name": "Updated name",
    "content": "Updated content"
  }
  ```
- `DELETE /api/v1/items/:id` - Delete an item

## Adding New Modules

To add a new module:

1. Create a new package in the `pkg/` directory
2. Create a module structure similar to the items module
3. Implement the module interface with a dependency on the web module
4. Register the module in the server.go file
5. The module will automatically be initialized when the server starts 