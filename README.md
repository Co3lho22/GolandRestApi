# GolangRestApi
## Overview

GolangRestApi is a robust and scalable RESTful API template built with Go (Golang). It's designed to serve as a solid foundation for developing various types of web applications. This API template includes features like user authentication, role-based access control (RBAC), token management, and more. It's structured to be easily extendable and customizable to fit the needs of different projects.
Features

**User Authentication:** Secure login and registration system.
    
**Token Management:** JWT-based authentication for secure API access.
    
**Role-Based Access Control (RBAC):** Fine-grained access control with roles and permissions.
    
**Admin Endpoints:** Specialized endpoints for administrative tasks.
    
**Middleware Integration:** Middleware for authentication and other common functionalities.

**Error Handling:** Standardized error responses for consistency and ease of debugging.

**Database Integration:** Ready-to-use database setup with MariaDB.

## Database Structure

This project uses a MariaDB database. To set up the database structure that this REST server supports, use the script provided in this repository. You can build and modify the database schema as needed for your specific application requirements.

## Getting Started

 To get started with GolangRestApi, follow these steps:

1. **Clone the Repository:** Clone this repository to your local machine.


    git clone https://github.com/Co3lho22/GolandRestApi.git

2. **Set Up the Database:** Use the script in this repository to set up your MariaDB database.


    git clone https://github.com/Co3lho22/basicRestAPIMariaDB/tree/main
    
3. **Configure the Application:** Create the .env following the structure of the .default_env.


4. **Build and Run:** Compile the application and start the server.


    go mod tidy

    go build -o GolandRestApi ./cmd/server

    ./GolandRestApi

## API Endpoints

The API includes the following endpoints:

    /api/v1/user/login:/ User login
    /api/v1/user/logout/{userId}: User logout
    /api/v1/user/register: User registration
    /api/v1/token/refresh: Token refresh
    /api/v1/admin/addUser: Add a new user (Admin only)
    /api/v1/admin/removeUser/{userId}: Remove a user (Admin only)

## Contributing

Contributions to improve GolangRestApi are welcome. Please feel free to submit pull requests or open issues to discuss proposed changes or enhancements.

## License

This project is licensed under the **GNU General Public License (GPL)**. This license allows users to freely use, modify, and distribute the software. However, if they distribute modified versions, they must also distribute the source code of their modifications under the GPL. This ensures that any modifications made to the software (if distributed) remain open-source, but it doesn't prevent modifications in forks.