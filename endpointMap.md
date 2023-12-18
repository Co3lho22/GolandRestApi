# API Endpoints Documentation

## User Login

        Endpoint: /login 
        Method: POST

Example Request (json):


    POST /login
    Authorization: Bearer <JWT Token>
    Content-Type: application/json
    
    {
    "username": "john_doe",
    "password": "password123"
    }

Example Response (json):

    HTTP/1.1 200 OK
    Content-Type: application/json

    {
      "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }

## User Registration

    Endpoint: /register
    Method: POST

Example Request (json):


    POST /register
    Authorization: Bearer <JWT Token>
    Content-Type: application/json
    
    {
    "username": "new_user",
    "password": "newpassword",
    "email": "newuser@example.com",
    "phone": "21323142134",
    "country": "Portugal"
    }

Example Response (json):


    HTTP/1.1 200 OK
    Content-Type: application/json

    "Registration Successful"

## Refresh Token

    Endpoint: /refresh
    Method: POST

Example Request (json):

    POST /refresh
    Authorization: Bearer <JWT Token>
    Content-Type: application/json
    
    {
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }

Example Response (json):

    HTTP/1.1 200 OK
    Content-Type: application/json

    {
      "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }

# Admin Operations

## Add User

    Endpoint: /admin/addUser
    Method: POST
    Authorization Required: Yes

Note: roleName can be "admin" or "user"

Example Request:

    POST /admin/addUser
    Authorization: Bearer <JWT Token>
    Content-Type: application/json

    {
    "username": "test22",
    "password": "test22",
    "email": "newuser@example.com",
    "country": "CountryName",
    "phone": "123456789",
    "role": {
                "name": "user"
            }
    }

Example Response:

    HTTP/1.1 200 OK
    Content-Type: application/json
    
    "User added successfull"

## Remove User

    Endpoint: /admin/removeUser/{userId}
    Method: DELETE
    Authorization Required: Yes

Example Request:

    DELETE /admin/removeUser/456
    Authorization: Bearer <JWT Token>

Example Response:

    HTTP/1.1 200 OK
    Content-Type: application/json
    
    "User removed successfully"