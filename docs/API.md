# API Documentation

## Overview

This API follows RESTful principles and uses JSON for request/response bodies.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API does not require authentication. Future versions will include JWT-based authentication.

## Response Format

All API responses follow this standard format:

```json
{
  "status": "success|error",
  "message": "Optional message",
  "data": "Response data (optional)",
  "timestamp": "2023-01-01T00:00:00Z"
}
```

## Endpoints

### Health Check

**GET** `/health`

Returns the health status of the service.

**Response:**
```json
{
  "status": "success",
  "message": "Service is healthy",
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### API Root

**GET** `/api/v1/`

Returns API information.

**Response:**
```json
{
  "status": "success",
  "message": "Clean Architecture API",
  "data": {
    "version": "1.0.0",
    "docs": "/docs"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### Users

#### List Users

**GET** `/api/v1/users`

Retrieves a list of users with pagination.

**Query Parameters:**
- `limit` (optional): Number of users to return (default: 10)
- `offset` (optional): Number of users to skip (default: 0)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "user_1234567890",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### Create User

**POST** `/api/v1/users`

Creates a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "user_1234567890",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### Get User

**GET** `/api/v1/users/{id}`

Retrieves a specific user by ID.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "user_1234567890",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### Update User

**PUT** `/api/v1/users/{id}`

Updates a specific user.

**Request Body:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": "user_1234567890",
    "email": "jane@example.com",
    "name": "Jane Doe",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### Delete User

**DELETE** `/api/v1/users/{id}`

Deletes a specific user.

**Response:**
```json
{
  "status": "success",
  "message": "User deleted successfully",
  "timestamp": "2023-01-01T00:00:00Z"
}
```

## Error Responses

When an error occurs, the API returns an error response:

```json
{
  "status": "error",
  "message": "Error description",
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## Rate Limiting

Currently, there is no rate limiting implemented. Future versions will include rate limiting.

## CORS

The API supports CORS and allows requests from any origin for development purposes. 