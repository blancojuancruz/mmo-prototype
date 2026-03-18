## Prerequisites

Before starting, ensure you have the following installed and configured:

* **Docker Desktop**: Installed and running.
* **Go (Golang)**: Make sure Go is installed on your system to run the server commands.

---

## Setup Instructions

Follow these steps to get the server up and running:

### 1. Start Docker
* Open **Docker Desktop** and wait until it is fully initialized and running.

### 2. Start the Database Container
* Always launch the Docker container first to ensure the database is accessible:
    ```bash
    docker-compose up -d
    ```

### 3. Launch the Server
* First, navigate to the server directory:
    ```bash
    cd server
    ```
* Then, run the server application:
    ```bash
    go run cmd/server/main.go
    ```