# Analytics Service API

## Overview
The Analytics Service processes events and provides insights about users activity. Most interactions happen asynchronously via Kafka events.

## Events Processed

### Tweet Created

**Topic**: `tweet_created`

**Schema**:
```json
{
  "id": "string",
  "event_type": "tweet_created",
  "handler": "string",
  "tweet_id": "string",
  "timestamp": "2025-08-09T05:13:41Z"
}
```

### Timeline Viewed

**Topic**: `timeline_viewed`

**Schema**:
```json
{
  "id": "string",
  "event_type": "timeline_viewed",
  "handler": "string",
  "timestamp": "2025-08-09T05:13:41Z"
}
```

## Endpoints

### Get User Analytics

```http
GET /v1/analytics/users/{id}
```

**Path Parameters**
- `id` (required): ID of the user

**Response**
```json
{
  "handler": "string",
  "is_influencer": true,
  "is_active": true,
}
```

### Get Users Analytics

```http
GET /v1/analytics/users
```

**Response**
```json
[
  {
    "handler": "string",
    "is_influencer": true,
    "is_active": true,
  }
]
```

### Delete User Analytics

```http
DELETE /v1/analytics/users/{id}
```

**Path Parameters**
- `id` (required): ID of the user

**Response**
```
204 No Content
```
