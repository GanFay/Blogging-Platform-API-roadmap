# 🚀 SecureBlog API

A **production-style REST API** for a blogging platform built with **Go**, **Gin**, and **PostgreSQL**.

The project implements a **secure authentication system** using **JWT access tokens**, **refresh tokens**, and **ownership protection** for blog posts.

It also includes:

- 🐳 Docker setup
- 📚 Swagger documentation
- 🔐 Secure authentication
- 📄 Pagination and filtering
- 🧱 Clean modular architecture

This project was built as a **backend portfolio project** demonstrating how a real-world API can be structured and implemented.

---

# ✨ Features

## 🔐 Authentication
- User registration
- User login
- JWT access tokens
- Refresh token flow
- Logout functionality
- `/users/me` endpoint to retrieve the current authenticated user

---

## 🛡 Security
- Password hashing using **bcrypt**
- JWT authentication middleware
- Protected routes
- Ownership validation (only the author can edit/delete their posts)

---

## 📝 Blog Posts
- Create post
- Update post
- Delete post
- Get all posts
- Get post by ID
- Search posts by text
- Pagination support (`limit` / `offset`)

---

## ⚙ API Infrastructure
- PostgreSQL database
- Docker Compose environment
- Environment variables configuration
- Swagger (OpenAPI) documentation
- Clean modular project architecture

---

# 🧰 Tech Stack

| Technology | Purpose |
|------------|--------|
| **Go** | Backend language |
| **Gin** | HTTP framework |
| **PostgreSQL** | Database |
| **pgxpool** | PostgreSQL driver |
| **JWT** | Authentication |
| **bcrypt** | Password hashing |
| **Docker** | Containerization |
| **Swagger (OpenAPI)** | API documentation |

---

# 📂 Project Structure

```
SecureBlog-API
│
├── auth/                 # JWT logic and password hashing
│   ├── password.go
│   └── token.go
│
├── docs/                 # Swagger documentation (generated)
│
├── handlers/             # HTTP handlers
│   ├── auth.go
│   ├── posts.go
│   ├── middleware.go
│   ├── me.go
│   └── ping.go
│
├── models/               # Data models
│   ├── post.go
│   └── user.go
│
├── router/               # Router configuration
│   └── router.go
│
├── postgres_data/        # PostgreSQL volume
│
├── Dockerfile
├── docker-compose.yml
├── .env
├── go.mod
├── main.go
└── README.md
```

---

# 🚀 Running the Project

## 🐳 Using Docker (recommended)

Start the API and PostgreSQL database:

```bash
docker compose up --build
```

API will be available at:

```
http://localhost:8080
```

Swagger documentation:

```
http://localhost:8080/swagger/index.html
```

![swagger.img](https://img.mtechlab.dev/uploads/e48ea0cdb2.webp)

---

## 💻 Running Locally (without Docker)

Install dependencies:

```bash
go mod tidy
```

Run the server:

```bash
go run main.go
```

---

# ⚙ Environment Variables

Create a `.env` file in the project root.

Example:

```
PG_USER=bloguser
PG_PASSWORD=admin
PG_DB=blogdb

JWT_SECRET=super_secret_jwt_key

APP_PORT=8080

DB_URL=postgres://bloguser:admin@localhost:5432/blogdb?sslmode=disable
```

---

# 📡 API Endpoints

## 🌐 Public Endpoints

| Method | Endpoint         | Description          |
| ------ | ---------------- | -------------------- |
| GET    | `/ping`          | Check server status  |
| POST   | `/auth/register` | Register a new user  |
| POST   | `/auth/login`    | Login user           |
| GET    | `/auth/refresh`  | Refresh access token |

---

## 🔒 Authenticated Endpoints

Require header:

```
Authorization: Bearer <access_token>
```

| Method | Endpoint | Description |
|------|------|------|
|GET | `/users/me` | Get current user |
|POST | `/auth/logout` | Logout user |
|POST | `/posts` | Create post |
|GET | `/posts` | Get all posts |
|GET | `/posts/:id` | Get post by ID |
|PUT | `/posts/:id` | Update post |
|DELETE | `/posts/:id` | Delete post |


---

# 📄 Pagination

The posts endpoint supports pagination.

Example:

```
GET /posts?limit=10&offset=0
```

Parameters:

| Parameter | Description |
|----------|-------------|
|limit | number of posts returned |
|offset | number of skipped posts |

Example:

```
GET /posts?term=golang&limit=5&offset=10
```

---

# 🔑 Authentication Flow

## Login

```
POST /auth/login
```

Response:

```json
{
  "access_token": "JWT_TOKEN"
}
```

---

## Access Protected Routes

Requests must include:

```
Authorization: Bearer <token>
```

---

## Refresh Token

```
GET /auth/refresh
```

Generates a new access token.

---

## Logout

```
POST /auth/logout
```

Removes refresh token.

---

# 🗄 Database Schema

## Users Table

```
users
```

| Column | Description |
|------|------|
|id | user id |
|username | username |
|email | user email |
|password_hash | hashed password |
created_at | account creation date |

---

## Posts Table

```
posts
```

| Column     | Description      |
| ---------- | ---------------- |
| id         | post id          |
| author_id  | post author      |
| title      | post title       |
| content    | post content     |
| category   | post category    |
| tags       | tags             |
| created_at | created time     |
| updated_at | last update time |


---

# 🛡 Security Features

- bcrypt password hashing
- JWT authentication
- token expiration
- refresh token flow
- ownership validation
- protected routes middleware

---

# 📚 Swagger Documentation

Interactive API documentation:

```
http://localhost:8080/swagger/index.html
```

Swagger allows you to:

- view all endpoints
- inspect request schemas
- test API directly in browser
- authenticate using JWT

---

# 🔄 Example Workflow

### Register

```
POST /auth/register
```

---

### Login

```
POST /auth/login
```

---

### Create Post

```
POST /posts
```

---

### Get Posts

```
GET /posts
```

---

# 🎯 Purpose of the Project

This project demonstrates how to build a **secure REST API backend** with:

- authentication
- database integration
- middleware
- pagination
- Docker infrastructure
- API documentation

It can serve as:

- 💼 a backend portfolio project
- 🚀 a starting point for a blogging platform
- 📚 a learning project for Go backend development

---

# 🔮 Future Improvements

Possible extensions:

- 💬 comments system
- ❤️ likes system
- 👥 role-based access control
- 🔍 full-text search
- 🚦 rate limiting
- 📦 database migrations
- ⚙ CI/CD pipeline
- ⚡ Redis caching

---

# 📜 License

MIT License
