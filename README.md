# Modular Go Application

A modular Go application with a web module using Gin and an items module that implements a REST API, protected by Clerk authentication.

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
│   ├── auth/             # Authentication module
│   │   ├── middleware.go # Clerk authentication middleware
│   │   └── module.go     # Auth module definition
│   ├── frontend/         # Frontend module (web UI)
│   │   └── module.go     # Frontend routes and handlers 
│   ├── items/            # Items module (items REST API)
│   │   ├── module.go     # Module definition and route handlers
│   │   ├── models.go     # Data models
│   │   └── repository.go # Data repository
│   └── web/              # Web module (core HTTP functionality)
│       └── module.go     # Module definition and base routing
├── web/
│   ├── static/           # Static assets (CSS, JS)
│   └── templates/        # HTML templates
├── go.mod                # Go module file
├── .env.example          # Example environment variables
└── README.md             # This file
```

## Modular Architecture

This application follows a modular architecture pattern:

- **Web Module**: Core HTTP server and routing functionality using Gin
- **Auth Module**: Authentication using Clerk
- **Items Module**: Business logic for item management, protected by authentication
- **Frontend Module**: Web UI for user interaction and authentication
- **Server**: Orchestrates and wires up all modules
- **Main**: Application entry point that starts the server

## Prerequisites

- Go 1.21 or later
- Clerk account for authentication (https://clerk.dev)

## Getting Started

1. Clone the repository
2. Create a `.env` file from `.env.example` and add your Clerk API keys:
   ```
   CLERK_API_KEY=your_clerk_api_key
   CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
   ```

3. Install dependencies:
   ```
   go mod download
   ```

4. Run the application:
   ```
   go run cmd/api/main.go
   ```

5. Access the application in your browser:
   ```
   http://localhost:8080/
   ```

## Using the Application

### Web Interface

1. Go to `http://localhost:8080/` in your browser
2. Click "Sign In" to authenticate using Clerk
3. After signing in, you'll be redirected to the dashboard
4. The dashboard displays your profile and API access information
5. You can copy your JWT token to use with the API

### Authentication

The application uses Clerk for authentication. All item endpoints are protected and require a valid JWT token.

To access protected endpoints:

1. Sign in through the web interface
2. Copy your JWT token from the dashboard
3. Include the token in your API requests:
   ```
   Authorization: Bearer your_jwt_token_here
   ```

## REST API Endpoints

The application provides a RESTful API for managing items, protected by Clerk authentication:

### Items Endpoints (All Require Authentication)

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

### Public Endpoints

- `GET /api/health` - Health check endpoint
- `GET /api/v1/ping` - Ping endpoint

## Testing with cURL

Testing protected endpoints:

```bash
# With a valid token
curl -H "Authorization: Bearer your_jwt_token_here" http://localhost:8080/api/v1/items

# Create a new item
curl -X POST -H "Authorization: Bearer your_jwt_token_here" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","content":"Content"}' \
  http://localhost:8080/api/v1/items
```

## Adding New Modules

To add a new module:

1. Create a new package in the `pkg/` directory
2. Create a module structure similar to the items module
3. Implement the module interface with dependencies on other modules as needed
4. Register the module in the server.go file
5. The module will automatically be initialized when the server starts 