# ChefBook Server
**ChefBook Server** is a proprietary REST API server for ChefBook App

## Architecture
* Language: Go
* DB: PostgreSQL
* Proxy-Server: Traefik

## Requirements
* Docker & Docker Compose

## Installation & Configuration
1. Clone the repo
2. Create and configure `.env` file in root directory:
```
# APP CONFIGURATION
APP_ENV=debug

# HTTP CONFIGURATION
HTTP_HOST=
HTTP_PORT=80
HTTPS_PORT=443

# FIREBASE CONFIGURATION
FIREBASE_API_KEY=

# DB CONFIGURATION
DB_NAME=
DB_PORT=
DB_USER=
DB_PASSWORD=

# BACKEND CONFIGURATION
BACKEND_PORT=
JWT_SIGNING_KEY=
SALT_COST=10

#S3 CONFIGURATION
S3_ACCESS_KEY=
S3_SECRET_KEY=

# SMTP CONFIGURATION
SMTP_EMAIL=
SMTP_PASSWORD=
```
3. Use `sudo docker-compose up` command to run server
