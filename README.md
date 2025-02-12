# DailyDiet Backend

This is the backend service for the DailyDiet application, a fullstack app designed to help users manage their daily diet and meal plans. The backend is built using Go and provides RESTful APIs for user authentication, meal management, and user statistics.

## Table of Contents

- [DailyDiet Backend](#dailydiet-backend)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Technologies](#technologies)
  - [Setup](#setup)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
  - [API Endpoints](#api-endpoints)
    - [Authentication](#authentication)
    - [Meals](#meals)
    - [User Statistics](#user-statistics)
  - [Contributing](#contributing)
  - [License](#license)

## Features

- User authentication (registration, login)
- Meal management (create, edit, delete meals)
- User statistics tracking (meals in diet, streaks)
- CORS support for cross-origin requests

## Technologies

- **Go**: Main programming language
- **Gin**: HTTP web framework
- **GORM**: ORM library for database interactions
- **PostgreSQL**: Database
- **Docker**: Containerization
- **godotenv**: Environment variable management
- **Testify**: Testing framework

## Setup

### Prerequisites

- Go 1.16+
- Docker
- PostgreSQL

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/aq-simei/DailyDietBackend.git
   cd daily_diet_backend
   ```

2. Create a [.env] file with the following content:

   ```env
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=daily_diet
   DB_PORT=5432
   DB_HOST=localhost
   JWT_SECRET=your_jwt_secret
   ```

3. Start PostgreSQL using Docker:

   ```bash
   docker-compose up -d
   ```

4. Install Go dependencies:

   ```bash
   go mod tidy
   ```

5. Run database migrations and seed data:
   ```bash
   go run main.go
   ```

## Usage

1. Start the server:

   ```bash
   go run main.go
   ```

2. The server will be running at `http://localhost:8080`.

## API Endpoints

### Authentication

- `POST /auth/register`: Register a new user
- `POST /auth/login`: Login a user

### Meals

- `POST /meals/new`: Create a new meal
- `GET /meals/list`: List all meals
- `PATCH /meals/edit/:mealId`: Edit a meal
- `DELETE /meals/delete/:mealId`: Delete a meal

### User Statistics

- `GET /user/stats`: Get user statistics

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Open a pull request

## License

This project is licensed under the MIT License.
