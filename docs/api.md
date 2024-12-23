# API Documentation

## Authentication

### Register User
```bash
curl -X POST http://localhost/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "yourpassword"
  }'
```

Response (201 Created):
```json
{
    "message": "User registered successfully"
}
```

### Login
```bash
curl -X POST http://localhost/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "yourpassword"
  }'
```

Response (200 OK):
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## Audio Processing

### Submit Audio Recording
POST `/audio/user/:user_id/phrase/:phrase_id`

Example:
```bash
curl -X POST http://localhost/audio/user/1/phrase/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: multipart/form-data" \
  -F "file=@/path/to/your/audio.m4a"
```

Response (201 Created):
```json
{
    "message": "Audio submitted successfully"
}
```

### Retrieve Audio Recording
GET `/audio/user/:user_id/phrase/:phrase_id/:format`

Example:

```bash
# Download WAV format
curl -X GET http://localhost/audio/user/1/phrase/1/wav \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  --output audio.wav

# Download M4A format
curl -X GET http://localhost/audio/user/1/phrase/1/m4a \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  --output audio.m4a
```

Response (200 OK):
- Content-Type: audio/wav or audio/m4a
- Binary audio file data

Supported formats:
- `wav` - WAV format
- `m4a` - M4A format

## Error Responses

All endpoints may return the following error responses:

```json
{
    "error": "Error message description"
}
```

Common HTTP status codes:
- 400 Bad Request - Invalid input
- 401 Unauthorized - Missing or invalid token
- 403 Forbidden - Insufficient permissions, e.g. accessing other user's audio
- 404 Not Found - Resource not found
- 500 Internal Server Error - Server error 