# Short - URL Shortener Service

A modern, efficient URL shortener service built with Go and featuring a clean web interface. This application allows users to create short, memorable links from long URLs while providing analytics and tracking capabilities.

## Features

- **URL Shortening**: Convert long URLs into short, manageable links
- **Analytics Tracking**: Monitor link usage with detailed analytics including:
  - Visitor country (using Cloudflare headers)
  - Referrer tracking
  - IP address logging
  - Timestamp tracking
- **Performance Optimized**:
  - LRU caching for frequently accessed URLs
  - Efficient database queries
  - Fast redirect handling
- **Security Features**:
  - Admin authentication system
  - Secure token handling
  - Environment-based configuration

## Technical Stack

- **Backend**: Go (Golang)
- **Database**: SQL database for persistent storage
- **Caching**: HashiCorp's LRU cache for performance optimization
- **Infrastructure**: Docker support for easy deployment
- **Development**: Hot-reload support with Air

## Getting Started

1. Clone the repository
2. Copy `example.env` to `.env` and configure your environment variables
3. Run the application:
   ```bash
   go run cmd/short/main.go
   ```

## Development

The project uses Air for hot-reloading during development. Start the development server with:
```bash
air
```

## Building

Build the application using Docker:
```bash
docker build -t short .
```

## License

This project is open source and available under the MIT License.
