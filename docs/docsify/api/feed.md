# Feed Service API

## Overview
The Feed Service provides the user timeline.

## Authentication
All endpoints require X-User-Id header.

## Endpoints

### Get User Timeline

```http
GET /timeline
```

**Query Parameters**
- `limit` (optional, default: 20): Number of tweets to return
- `offset` (optional, default: 0): Pagination offset

**Headers**
- `X-User-Id` (required): ID of the user

**Response**
```json
{
  "tweets": [
    {
      "id": "string",
      "handler": "string",
      "content": {
        "text": "string",
      },
      "created_at": "2025-08-09T05:13:41Z"
    }
  ],
  "next_offset": 20
}
```
