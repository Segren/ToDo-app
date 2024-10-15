# ToDo App with REST API

This is a simple ToDo application built with Go (Golang), Gin framework, GORM, PostgreSQL, and JWT authentication. The application allows users to register, log in, and manage their tasks (create, read, update, and delete tasks). It also includes Swagger documentation for the API.


## Features
- User authentication with JWT (JSON Web Tokens)
- Create, read, update, and delete tasks (CRUD operations)
- PostgreSQL database integration
- Swagger documentation for the API
- Docker support for easy deployment


## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Docker Setup](#docker-setup)
- [Testing](#testing)


## Installation

### Prerequisites
- Go 1.23 or higher
- PostgreSQL 13 or higher
- Git

### Step 1: Clone the repository
bash
git clone https://github.com/<your-username>/todo-app.git
cd todo-app

### Step 2: Clone the repository

Create a .env file in the root directory of the project to store the database credentials and other configurations:

DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todoapp
DB_PORT=5432

### Step 3: Install dependencies

go mod tidy

### Step 4: Set up PostgreSQL

Make sure you have PostgreSQL installed and running. You can create the todoapp database with the following command:

```createdb todoapp```

### Step 5: Run the application

from cmd/toDo-app/ run
``` go run main.go ```

The server will be running on http://localhost:8080.


## Usage

### Registration

To register a new user, send a POST request to /register with the following JSON body:

```
{
  "username": "testuser",
  "password": "password123"
}
```

### Login

To log in, send a POST request to /login with the following JSON body:

```
{
  "username": "testuser",
  "password": "password123"
}
```

You will receive a JWT token in response. Use this token to access protected routes by including it in the Authorization header:

```Authorization: Bearer <your-token>```

### Task Management

Once logged in, you can manage your tasks (CRUD operations) with the following endpoints:

	•	GET /todos: Get all tasks
	•	POST /todos: Create a new task
	•	PUT /todos/:id: Update an existing task
	•	DELETE /todos/:id: Delete a task

### API Documentation

The API is documented using Swagger. Once the application is running, you can access the documentation by navigating to:

http://localhost:8080/swagger/index.html


## Docker Setup

This project includes Docker support for easier setup and deployment.

### Build and run with Docker Compose

```cd docker/```

```docker-compose up --build```

This will set up the application and a PostgreSQL database. The application will be available at http://localhost:8080.
