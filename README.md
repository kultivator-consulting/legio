# Cortex API

## Overview
Cortex API is a microservices-based application built with Go. It includes several services such as authentication, file management, content management, and more. The project uses Docker for containerization and PostgreSQL as the database.

## Services
- **auth_service**: Handles authentication and authorization.
- **file_service**: Manages file storage and retrieval.
- **cms_service**: Content management service.
- **cortex_admin**: Admin interface for managing the application.
- **cortex_website**: Frontend website for the application.

## Configuration
### `services/auth_service/config.yaml`
```yaml
bundleId: cloud.legio.auth_api
cookieDomain: "legio.cloud"
accessTokenExpiresIn: 60m
accessTokenMaxAge: 60m
refreshTokenExpiresIn: 24h
refreshTokenMaxAge: 24h
```

## Docker Setup
The project uses Docker Compose to manage the services. Below is a brief description of each service defined in `docker-compose.yaml`:

- **db**: PostgreSQL database.
- **nginx**: Nginx server for handling HTTP requests.
- **migrate**: Database migration service.
- **test**: Service for running tests.
- **auth_service**: Authentication service.
- **file_service**: File management service.
- **cms_service**: Content management service.
- **cortex_admin**: Admin interface.
- **cortex_website**: Frontend website.

### Running the Application
To start the application, run:
```sh
docker-compose up
```

### Stopping the Application
To stop the application, run:
```sh
docker-compose down
```

## Testing
The `test` service runs Go tests and generates a JUnit report. To execute the tests, use:
```sh
docker-compose run test
```

## Network
All services are connected through a Docker network named `cortex-net`.

## License
This project is licensed under the MIT License.
```

This `README.md` provides an overview of the project, configuration details, Docker setup, and instructions for running and testing the application.