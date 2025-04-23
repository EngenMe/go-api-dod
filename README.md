# Go API with DOD Principles

This is a backend API boilerplate built with Go, following Data-Oriented Design (DOD) principles. It uses GORM for PostgreSQL database operations and Gin for the HTTP framework.

## Features

- Complete user management system (CRUD operations)
- Authentication system with JWT tokens
- Password hashing with bcrypt
- PostgreSQL database integration with GORM
- Structured error handling
- Environment-based configuration
- Request logging

## Data-Oriented Design Principles

This project follows DOD principles:
- Procedural code with clear data types
- Data flow optimization
- Bundling global variables into structs
- Module-level encapsulation
- Parameterized functions
- Focused on data organization and access patterns

## Folder Structure

```
go-api-dod/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration loader
├── internal/
│   ├── api/
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWT authentication middleware
│   │   │   └── logging.go          # Request logging middleware
│   │   ├── handlers/
│   │   │   ├── auth.go             # Authentication handlers
│   │   │   └── users.go            # User CRUD handlers
│   │   └── server.go               # API server setup
│   ├── data/
│   │   ├── models/
│   │   │   └── user.go             # User model definition
│   │   └── store/
│   │       ├── postgres.go         # Database connection
│   │       └── users.go            # User data operations
│   └── utils/
│       ├── hash.go                 # Password hashing utilities
│       └── token.go                # JWT token utilities
├── migrations/
│   └── 001_create_users_table.sql  # Database migration
├── .env.example                    # Example environment variables
├── go.mod                          # Go module definition
└── README.md                       # Project documentation
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

## Setup and Usage

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-api-dod.git
   cd go-api-dod
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up environment variables:
   ```
   cp .env.example .env
   ```
   Edit the `.env` file to match your environment.

4. Set up the database:
   ```
   createdb goapi
   ```

5. Build and run the server:
   ```
   go build -o api ./cmd/api
   ./api
   ```
   Or simply run:
   ```
   go run ./cmd/api
   ```

## API Endpoints

### Authentication

- `POST /signup` - Register a new user
    - Request: `{ "email": "user@example.com", "password": "password123" }`
    - Response: `{ "token": "JWT_TOKEN", "user": { "id": "UUID", "email": "user@example.com" } }`

- `POST /login` - Login with existing user
    - Request: `{ "email": "user@example.com", "password": "password123" }`
    - Response: `{ "token": "JWT_TOKEN", "user": { "id": "UUID", "email": "user@example.com" } }`

### Users (Protected Routes - Requires Authorization Header)

- `GET /users` - List all users
    - Headers: `Authorization: Bearer JWT_TOKEN`
    - Query Parameters: `limit=10&offset=0`
    - Response: `[{ "id": "UUID", "email": "user@example.com", "created_at": "TIMESTAMP", "updated_at": "TIMESTAMP" }]`

- `GET /users/:id` - Get a user by ID
    - Headers: `Authorization: Bearer JWT_TOKEN`
    - Response: `{ "id": "UUID", "email": "user@example.com", "created_at": "TIMESTAMP", "updated_at": "TIMESTAMP" }`

- `POST /users` - Create a new user
    - Headers: `Authorization: Bearer JWT_TOKEN`
    - Request: `{ "email": "newuser@example.com", "password": "password123" }`
    - Response: `{ "id": "UUID", "email": "newuser@example.com", "created_at": "TIMESTAMP" }`

- `PUT /users/:id` - Update a user
    - Headers: `Authorization: Bearer JWT_TOKEN`
    - Request: `{ "email": "updated@example.com" }`
    - Response: `{ "id": "UUID", "email": "updated@example.com", "created_at": "TIMESTAMP", "updated_at": "TIMESTAMP" }`

- `DELETE /users/:id` - Delete a user
    - Headers: `Authorization: Bearer JWT_TOKEN`
    - Response: `{ "message": "User deleted successfully" }`

## License

This project is licensed under the MIT License - see the LICENSE file for details.