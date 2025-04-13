# Fitness Platform API

A production-ready backend for a fitness platform built with Go, Fiber, and PostgreSQL.

## 🚀 Features

### Authentication & Authorization
- OTP-based authentication for mobile users
- Email token-based authentication for trainers/admins
- Role-based access control (Customer, Trainer, Gym Owner)
- Secure session management with JWT tokens

### User Management
- Customer registration and profile management
- Trainer registration and profile management
- Gym owner registration and profile management
- User role management and permissions

### Gym Management
- Gym registration and profile management
- Gym location and contact information
- Gym facilities and amenities tracking
- Gym operating hours management

### Training & Workout Management
- Workout plan creation and management
- Exercise tracking and logging
- Training session scheduling
- Progress tracking and analytics

### Booking & Scheduling
- Training session booking
- Class scheduling
- Appointment management
- Availability management for trainers

### Payment & Subscription
- Subscription plan management
- Payment processing
- Billing and invoicing
- Membership management

### Communication
- In-app messaging
- Notification system
- Announcement broadcasting
- Feedback and review system

## 🛠️ Tech Stack

- Language: Go
- Web Framework: Fiber
- Database: PostgreSQL
- Cache/OTP store: Valkey (Redis compatible)
- ORM: GORM
- Authentication: OTP-based (mobile), Email token-based (for owners/trainers)
- API Documentation: Swagger
- Migration Tool: golang-migrate

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis/Valkey
- Make
- golang-migrate

## 🚀 Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/yourname/fitness-platform.git
   cd fitness-platform
   ```

2. Install dependencies:

   ```bash
   make deps
   ```

3. Copy the environment file:

   ```bash
   cp .env.example .env
   ```

4. Update the `.env` file with your configuration

5. Run database migrations:

   ```bash
   make migrate-up
   ```

6. Start the server:
   ```bash
   make run
   ```

The server will start at `http://localhost:8080`

## 📚 API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger
```

## 🛠️ Development

### Available Make Commands

- `make run` - Run the application
- `make build` - Build the application
- `make test` - Run tests
- `make clean` - Clean build files
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make swagger` - Generate Swagger documentation
- `make deps` - Install dependencies
- `make lint` - Run linter
- `make security-check` - Run security checks
- `make migration` - Create a new migration

### Project Structure

```
fitness-platform/
├── cmd/                    # Application entry points
├── config/                 # Configuration
├── internal/              # Private application code
├── pkg/                   # Public packages
├── migrations/            # Database migrations
├── scripts/              # Utility scripts
├── docs/                 # Documentation
└── tests/                # Integration tests
```

## 🔒 Security

- All passwords are hashed using bcrypt
- JWT tokens for session management
- Rate limiting on API endpoints
- CORS protection
- Input validation
- SQL injection protection through GORM
- XSS protection

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.
