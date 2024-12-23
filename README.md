# SBC

SBC is a simple backend service that can accept, convert, store, and retrieve an audio file associated with a user and a practice phrase.

## Requirements

- Docker and Docker Compose
- Go 1.21 or later (for local development)
- FFmpeg (for audio processing)

## Environment Setup

1. Create a `.env` file in the root directory with the following variables:
```
DB_USER=
DB_PASSWORD=
DB_ROOT_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
JWT_SECRET=
```
Or copy the `.env.example` file and edit the values.

2. Install docker and docker compose.

## Run the service

1. Start the containers:
```
docker compose up -d
```

2. The service will be available at `http://localhost`.

## Documentation

- [How Audio Processing Works](docs/audio-processing.md)
- [API Documentation](docs/api.md)

## Default Credentials

The service comes with a default user and 3 phrases:
- User:
  - Email: admin@example.com
  - Password: admin123
- Phrases:
  - PhraseID 1: "Hello"
  - PhraseID 2: "SpeakBuddy is an online service that helps users acquire language abilities"
  - PhraseID 3: "We believe that the spoken language sends a message not to the brain, but to the heart"

