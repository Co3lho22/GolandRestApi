# GolangRestApi
This is a general, but robust rest api in golang that can be used for future projects

# DB structure 
Use this script https://github.com/Co3lho22/basicRestAPIMariaDB to set up the DB that is supported by this rest server 
and then build on it if needed.

# Project Structure

my-golang-api/
│
├── cmd/
│   └── server/
│       └── main.go        # Entry point of the application
│
├── pkg/
│   ├── api/               # API-related logic (controllers, middleware)
│   │   ├── handlers/      # Request handlers
│   │   └── middleware/    # API middleware
│   │
│   ├── config/            # Configuration related logic
│   │   └── config.go      # Configuration struct and loader
│   │
│   ├── model/             # Data models
│   │   └── model.go       # Structs for data representation
│   │
│   ├── service/           # Business logic
│   │   └── service.go     # Service layer logic
│   │
│   └── repository/        # Data access layer
│       └── repository.go  # Database interactions
│
├── internal/              # Internal packages (not for external use)
│
├── test/                  # Test files
│
├── Dockerfile             # Docker configuration
│
├── .env                   # Environment variables
│
├── go.mod                 # Go module definitions
└── go.sum                 # Go module checksums
