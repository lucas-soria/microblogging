# Tweets Service API

## Overview
The Tweets Service provides the user tweets.

## Authentication
All endpoints require X-User-Id header.

## Endpoints

### Create Tweet

```http
POST /tweets
```

**Headers**
- `X-User-Id` (required): ID of the user

**Request Body**
```json
{
  "content": {
    "text": "Hello, world!"
  },
  "handler": "string"
}
```

**Response**
```json
{
  "id": "string",
  "handler": "string",
  "content": {
    "text": "string",
  },
  "created_at": "2025-08-09T05:13:41Z"
}
```

### Get Tweet

```http
GET /tweets/{id}
```

**Path Parameters**
- `id` (required): ID of the tweet

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```json
{
  "id": "string",
  "handler": "string",
  "content": {
    "text": "string",
  },
  "created_at": "2025-08-09T05:13:41Z"
}
```

### Get User Tweets

```http
GET /tweets/users/{id}
```

**Path Parameters**
- `id` (required): ID of the user

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```json
[
  {
    "id": "string",
    "handler": "string",
    "content": {
      "text": "string",
    },
    "created_at": "2025-08-09T05:13:41Z"
  }
]
```

### Delete Tweet

```http
DELETE /tweets/{id}
```

**Path Parameters**
- `id` (required): ID of the tweet to delete

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```
204 No Content
```
