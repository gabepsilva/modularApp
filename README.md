# Modular Go Application

A modular Go application with a web module using Gin, an items module that implements a REST API, and Clerk authentication with Memcached caching for JWT tokens.

## Project Structure

```
modularApp/
├── cmd/
│   └── api/              # Application entry points
│       └── main.go
├── dev_containers/
│   └── memcached.podman-compose.yml # Memcached containers configuration
├── internal/
│   └── server/           # Internal server implementation
│       └── server.go
├── pkg/
│   ├── auth/             # Authentication module
│   │   ├── middleware.go # Clerk authentication middleware
│   │   └── module.go     # Auth module definition
│   ├── cache/            # Caching module 
│   │   ├── module.go     # Cache module definition
│   │   ├── jwt_cache.go  # JWT token caching
│   │   └── memcached.go  # Memcached client
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
├── main.go               # Main application entry point
├── .env.example          # Example environment variables
└── README.md             # This file
```

## Modular Architecture

This application follows a modular architecture pattern:

- **Web Module**: Core HTTP server and routing functionality using Gin
- **Auth Module**: Authentication using Clerk with JWT token validation
- **Cache Module**: Caching layer using Memcached for JWT tokens
- **Items Module**: Business logic for item management, protected by authentication
- **Frontend Module**: Web UI for user interaction and authentication
- **Server**: Orchestrates and wires up all modules
- **Main**: Application entry point that starts the server

## Prerequisites

- Go 1.21 or later
- Clerk account for authentication (https://clerk.dev)
- Memcached servers (can be run using Podman or Docker)

## Getting Started

1. Clone the repository

2. Set up Memcached servers:
   ```
   cd dev_containers
   podman-compose -f memcached.podman-compose.yml up -d
   ```
   
   Or with Docker:
   ```
   cd dev_containers
   docker-compose -f memcached.podman-compose.yml up -d
   ```

3. Create a `.env` file from `.env.example` and add your Clerk API keys:
   ```
   # Clerk Authentication
   CLERK_API_KEY=your_clerk_api_key_here
   CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key_here

   # Memcached servers (comma-separated)
   MEMCACHED_SERVERS=localhost:11211,localhost:11212

   # JWT Cache TTL in seconds (default is 900 seconds for 15 minutes)
   JWT_CACHE_TTL_SECONDS=900
   ```

4. Install dependencies:
   ```
   go mod download
   ```

5. Run the application:
   ```
   go run main.go
   ```

6. Access the application in your browser:
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

The application uses Clerk for authentication with Memcached caching for JWT tokens:

1. JWT tokens are validated using Clerk's API
2. Validated tokens are stored in Memcached for faster subsequent requests
3. This improves performance by reducing API calls to Clerk
4. All item endpoints are protected and require a valid JWT token

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

## Caching Layer

The application uses Memcached for caching JWT tokens:

1. Two Memcached instances are used for redundancy
2. JWT tokens are cached with a default 15-minute TTL
3. Caching improves performance by reducing calls to Clerk API
4. The cache is load-balanced across available Memcached servers

## Adding New Modules

To add a new module:

1. Create a new package in the `pkg/` directory
2. Create a module structure similar to the items module
3. Implement the module interface with dependencies on other modules as needed
4. Register the module in the server.go file
5. The module will automatically be initialized when the server starts

## Development Notes

- The app uses a modular architecture with DI (Dependency Injection)
- All modules are initialized in `server.New()` 
- Environment variables are loaded in each module using direct os.Getenv calls
- Authentication middleware is created with JWT caching support
- Memcached is used in a load-balanced setup with multiple servers

## TODO and Improvement Ideas

The following are potential improvements and enhancements for the application:

1. **Configuration Management**:
   - Replace direct `os.Getenv()` calls with a centralized configuration package
   - Add support for configuration files (YAML/JSON) alongside environment variables
   - Implement validation for required configuration values

2. **Testing**:
   - Add unit tests for all modules with good coverage
   - Implement integration tests for API endpoints
   - Add mock implementations for external dependencies (Clerk, Memcached)

3. **Logging**:
   - Replace `log.Printf` with a structured logging package (e.g., zap or logrus)
   - Add request ID tracking for better request tracing
   - Implement different log levels based on environment (dev/prod)

4. **Error Handling**:
   - Create consistent error types and responses across the application
   - Add more detailed error messages for debugging
   - Implement proper error propagation from lower layers

5. **Security Enhancements**:
   - Add rate limiting for API endpoints
   - Implement CORS configuration for API
   - Add security headers to all responses
   - Consider adding a CSP (Content Security Policy)

6. **Performance**:
   - Implement database connection pooling
   - Add response compression
   - Consider implementing request timeouts

7. **Monitoring and Observability**:
   - Add Prometheus metrics
   - Implement health check endpoints with detailed component status
   - Add distributed tracing (OpenTelemetry/Jaeger)

8. **DevOps**:
   - Create Dockerfile and docker-compose for the entire application
   - Add CI/CD pipeline configuration
   - Implement graceful shutdown with proper signal handling

9. **Feature Additions**:
    - Add user management features
    - Implement role-based access control
    - Add pagination, filtering, and sorting for list endpoints
    - Create a more robust frontend with modern JavaScript framework 