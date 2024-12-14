# Go Lang User Management API

This is a Go project that uses Docker and the [Gin-Gonic](https://github.com/gin-gonic/gin) framework to build a simple REST API for user management. The API supports CRUD (Create, Read, Update, Delete) operations for users and sends a "Happy Birthday" email to users when it's their birthday.

## Features
- Create, Read, Update, and Delete user records.
- Send a "Happy Birthday" email to users on their birthday.

## Requirements
- Go 1.18+
- Docker
- Gin-Gonic
- MongoDB

## Setup

### Clone the Repository
Clone the repository to your local machine:

``` bash
git clone https://github.com/yourusername/projectname.git
cd projectname
 ```

Build and Run with Docker
Build the Docker image:
``` bash 
docker build -t user-api .
Run the container:
docker run -p 8080:8080 user-api
```

## Endpoints
- POST /users - Create a new user
- GET /users/:id - Get a user by ID
- PUT /users/:id - Update a user by ID
- DELETE /users/:id - Delete a user by ID
- POST /user/sendemail - Send email

## Email Notification
The API checks daily for users whose birthdays match the current date and sends them a "Happy Birthday" email.

## License
This project is licensed under the MIT License.
