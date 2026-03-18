## Prerequisites

Before starting, ensure you have the following installed and configured:

* **Docker Desktop**: Installed and running.
* **Go (Golang)**: Make sure Go is installed on your system to run the server commands.

---

## Setup Instructions

Follow these steps to get the server up and running:

### 1. Environment Configuration
* Copy the example environment file and configure your local variables:
    ```bash
    cp .env.example .env
    ```

### 2. Start Docker
* Open **Docker Desktop** and wait until it is fully initialized and running.

### 3. Start the Database Container
* Always launch the Docker container first to ensure the database is accessible:
    ```bash
    docker-compose up -d
    ```

### 4. Launch the Server
* First, navigate to the server directory:
    ```bash
    cd server
    ```
* Then, run the server application:
    ```bash
    go run cmd/server/main.go
    ```

---

## Troubleshooting

If you encounter issues during the setup, check the following:

* **Port Conflict:** If `docker-compose up` fails, ensure the database port is not being used by another local service.
* **Docker Daemon:** If you get a "cannot connect to the Docker daemon" error, make sure Docker Desktop is actually running.
* **Dependencies:** If `go run` fails, try running `go mod tidy` in the `server` folder to install any missing dependencies.