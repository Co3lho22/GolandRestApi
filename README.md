# GolangRestApi
## Overview

GolangRestApi is a robust and scalable RESTful API template built with Go (Golang). It's designed to serve as a solid foundation for developing various types of web applications. This API template includes features like user authentication, role-based access control (RBAC), token management, and more. It's structured to be easily extendable and customizable to fit the needs of different projects.

## Features

**User Authentication:** Secure login and registration system.
    
**Token Management:** JWT-based authentication for secure API access.
    
**Role-Based Access Control (RBAC):** Fine-grained access control with roles and permissions.
    
**Admin Endpoints:** Specialized endpoints for administrative tasks.
    
**Middleware Integration:** Middleware for authentication and other common functionalities.

**Error Handling:** Standardized error responses for consistency and ease of debugging.

**Database Integration:** Ready-to-use database setup with MariaDB.

## Database Structure

This project uses a MariaDB database. The `init-db.sql` script provided in this repository sets up the database structure, including a default admin account (**username: admin**, **password: admin**). You can build and modify the database schema as needed for your specific application requirements.

## Persistence with Docker Volumes

The **MariaDB** database uses a Docker volume to ensure **data persistence**. This means that your data remains intact even when the database container is stopped or restarted. The volume is defined in the `docker-compose.yml` file under the `volumes` section for the `db` service.

## Getting Started

 To get started with GolangRestApi, follow these steps:

1. **Clone the Repository:** 

    ```bash
    git clone https://github.com/Co3lho22/GolandRestApi.git
    ```
   
2. **Configure the Application:** Create the .env following the structure of the .default_env.


4. **Build and Run with Docker:** Compile the application and start the server.
    ```bash
    docker-compose up -d --build
    ```

## API Endpoints

The API includes the following endpoints:

* **/api/v1/user/login:** User login 

```bash
   curl -X POST http://localhost:8080/api/v1/user/login -d '{"username":"<username>", "password":"<password>"}'
```

* **/api/v1/user/logout/{userId}:** User logout 

```bash
   curl -X GET http://localhost:8080/api/v1/user/logout/{userId}
```

* **/api/v1/user/register:** User registration 

```bash
curl -X POST http://localhost:8080/api/v1/user/register -d '{"username":"<username>", "password":"<password>", "email":"<email>"}'
```

* **/api/v1/token/refresh:** Token refresh

```bash
curl -X POST http://localhost:8080/api/v1/token/refresh -d '{"refreshToken":"<refreshToken>"}'
```

* **/api/v1/admin/addUser:** Add a new user (Admin only)

```bash
curl -X POST http://localhost:8080/api/v1/admin/addUser -d '{"user": {"username":"<username>", "password":"<password>", "email":"<email>"}, "roleName":"<roleName>"}'
```

* **/api/v1/admin/removeUser/{userId}:** Remove a user (Admin only)

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/removeUser/{userId}
```
Replace **`<username>`**, **`<password>`**, **`<email>`**, **`<refreshToken>`**, **`<roleName>`**, and **`{userId}`** with appropriate values for your tests.

**Note:** you might need to adapt the url endpoint depending on your .env file configuration.

## Containerization

The application is containerized using Docker and managed with Docker Compose. This setup includes separate containers for the REST API server and the MariaDB database. The docker-compose.yml file simplifies deployment and ensures consistency across different environments.

## Docker Image Options

The `Dockerfile` in the repository is set up to support two different build strategies for the Docker image:

1. **Full Code Image:**
   * This build includes all the source code along with the compiled binary. It's useful for environments where you might want to inspect or modify the source code within the container.   
   * The relevant section of the `Dockerfile` for this build is:
   
      ```bash
      # Build stage
      FROM golang:1.21.5-alpine3.19 AS builder
      
      WORKDIR /restApi
      
      COPY go.mod go.sum ./
      RUN go mod download
      COPY . .
      RUN go build -o GolandRestApi ./cmd/server
      
      # Final stage - all the code
      EXPOSE 8080
      CMD ["./GolandRestApi"]
      ```
2. **Executable-Only Image:**
   * This build includes only the compiled executable in a minimal Alpine Linux environment. It's a lightweight option, ideal for production deployments where you don't need the source code.
   * To use this build, uncomment the following lines in the `Dockerfile`:
       
       ```bash
       # Final stage - only with the executable
       FROM alpine:3.19.0
   
       WORKDIR /root/
       COPY --from=builder /restApi/GolandRestApi .
       EXPOSE 8080
       CMD ["./GolandRestApi"]
       ```
Choose the build strategy that best fits your deployment needs. The full code image is recommended for development environments, while the executable-only image is more suited for production deployments.

## Interacting with Containers
1. **Accessing the MariaDB Container:**
   * To access the MariaDB database, use the following command:
        ```bash
        docker exec -it golandrestapi-db-1 mariadb -u restServer -p
        ```
   You will be prompted to enter the password for the restServer user that you defined in the .env file

2. **Accessing the REST API Container:**
   * To access the shell of the REST API container, use the following command:
        ```bash
        docker exec -it golandrestapi-restapi-1 sh
        ```

3. **Viewing Logs:**
   * If you want to view the logs of the REST API, you can use the following command: 
        ```bash
        docker logs -f golandrestapi-restapi-1
        ```

4. **Stopping Containers:**
    * To stop the running containers, you can use the following command:
        ```bash
        docker-compose down
        ```

5. **Rebuilding Containers:**
    * If you make changes to your application and need to rebuild the containers, use:
        ```bash
        docker-compose up -d --build
        ```
   This command rebuilds the containers with the latest changes.

6. **Listing Active Containers:**
    * To see a list of all active containers, use:
        ```bash
        docker ps
        ```


## Current Work-in-Progress and TODOs

* **Update Config to Use Docker Environment Variables:** Modify the application to directly use environment variables set in docker-compose.yml.
* **Volume for Logs:** Consider mounting the log directory (var) as a volume for persistent log storage.
* **Implement Docker Secrets:** Update the application to use Docker secrets for sensitive data instead of relying on .env files.



## Contributing

Contributions to improve GolangRestApi are welcome. Please feel free to submit pull requests or open issues to discuss proposed changes or enhancements.

## License

This project is licensed under the **GNU General Public License (GPL)**. This license allows users to freely use, modify, and distribute the software. However, if they distribute modified versions, they must also distribute the source code of their modifications under the GPL. This ensures that any modifications made to the software (if distributed) remain open-source, but it doesn't prevent modifications in forks.
