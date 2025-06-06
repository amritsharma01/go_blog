# 📝 Go Blog API

A RESTful blog API built with **Go**, using **Echo** as the web framework and **GORM** as the ORM. The application follows a clean layered architecture with well-separated concerns (handlers, services, repositories), JWT authentication, and role-based protected routes.

## 🚀 Features

- User registration and login with JWT authentication
- Secure protected routes for users and authors
- CRUD operations for blog posts
- Paginated post listing and author-specific views
- Category management (create, delete, list)
- Clean architecture: repository, service, and handler layers
- Swagger documentation for all API endpoints
- Custom error handling and middleware support

---

## 🛠️ Tech Stack

- **Go (Golang)**
- **Echo** – Web framework
- **GORM** – ORM for database operations
- **JWT** – Authentication and protected routes
- **Swagger** – Auto-generated API docs
- **PostgreSQL** – Primary database (can be swapped)

---

## 📁 Project Structure
```
go_blog/
├── cmd/
  └── main.go # Application entry point
├── handlers/ # HTTP layer (Echo handlers)
├── middleware/ # JWT and custom middleware
├── repositories/ # Data access layer
├── services/ # Business logic layer
├── models/ # GORM models
├── errors/ # Custom error handlers
├── request_models/ # Request DTOs
├── response_models/ # Response DTOs
├── routes/ # Route definitions
├── config/ # Configuration and environment setup
├── utils/ # Utility functions
├── docs/ # Swagger generated files
├── go.mod # Go module file
└── go.sum # Go module checksum
```

## 🔐 Authentication

This project uses **JWT** for stateless authentication.

- `POST /v1/auth/register` – Register a new user
- `POST /v1/auth/login` – Get JWT token
- Protected routes require `Authorization: Bearer <token>` header

---

## 📚 API Endpoints

### Public

- `GET /` – Welcome message
- `GET /v1/posts` – List all posts (paginated)
- `GET /v1/posts/:id` – Get post by ID
- `GET /swagger/*` – Swagger API documentation

### Auth

- `POST /v1/auth/register` – Register
- `POST /v1/auth/login` – Login

### Users (Protected)

- `GET /v1/users` – List all users

### Posts (Protected)

- `POST /v1/posts` – Create a new post
- `PATCH /v1/posts/:id` – Edit post
- `DELETE /v1/posts/:id` – Delete post
- `GET /v1/authors/:author_id/posts` – Get posts by author (paginated)

### Categories (Protected)

- `GET /v1/categories` – List all categories (paginated)
- `POST /v1/categories` – Create a category
- `DELETE /v1/categories/:id` – Delete a category

---


## 🧪 Running Locally

### 1. Clone the repo

```bash
git clone https://github.com/amritsharma01/go_blog.git
cd go_blog
```

### 2. Setup `.env`

Create a `.env` file in the root of your project and configure your database and JWT secret:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_blog

JWT_SECRET=your_jwt_secret_key
```
### 3. Run the project
```bash
go mod tidy
go run cmd/main.go
```
#### Now the project is accessible at *localhost:8000*

## Swagger Documentation
### Swagger Documentation is availabe at 
```
GET /swagger/index.html
```
### To regenerate the swagger files, run
```bash
swagger init
```
## Testing
### You can use Postman or cURL for testing the endpoints. Automatic tests for services and handlers are unser development.

## Contributors
- **Amrit Sharma** - [amritsharma1027@gmail.com](mailto:amritsharma1027@gmail.com)

## Contact
For questions or feedback, reach out via email.




