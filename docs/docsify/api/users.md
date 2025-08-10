# Users Service API

## Overview
The Users Service provides the user and follow management.

## Authentication
All endpoints require X-User-Id header.

## Endpoints

### Create User

```http
POST /users
```

**Headers**
- `X-User-Id` (required): ID of the user

**Request Body**
```json
{
  "first_name": "string",
  "last_name": "string",
  "handler": "string"
}
```

**Response**
```json
{
  "handler": "string",
  "first_name": "string",
  "last_name": "string"
}
```

### Get User

```http
GET /users/{id}
```

**Path Parameters**
- `id` (required): ID of the user

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```json
{
  "handler": "string",
  "first_name": "string",
  "last_name": "string"
}
```

### Delete User

```http
DELETE /users/{id}
```

**Path Parameters**
- `id` (required): ID of the user

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```
204 No Content
```

### Follow User

```http
POST /users/{id}/follow
```

**Path Parameters**
- `id` (required): ID of the user to follow

**Headers**
- `X-User-Id` (required): ID of the user

**Request Body**
```json
{
  "followee_id": "string"
}
```

**Response**
```
202 Accepted
```

### Unfollow User

```http
POST /users/{id}/unfollow
```

**Path Parameters**
- `id` (required): ID of the user to follow

**Headers**
- `X-User-Id` (required): ID of the user

**Request Body**
```json
{
  "followee_id": "string"
}
```

**Response**
```
204 No Content
```

### Get User Followees

```http
GET /users/{id}/followees
```

**Path Parameters**
- `id` (required): ID of the user

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```json
[
  {
    "handler": "string",
    "first_name": "string",
    "last_name": "string"
  }
]
```
