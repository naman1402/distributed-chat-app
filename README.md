## Distributed Chat App - Readme.md

This is a distributed chat application built using Golang for the backend, Redis for Pub/Sub messaging, websockets for real-time communication, Docker for containerization, Nginx for reverse proxying, and Cassandra for persistent data storage.

### Technologies Used

* **Golang:** Backend programming language
* **Redis:** Pub/Sub messaging server
* **Websockets:** Real-time communication protocol
* **Docker:** Containerization platform
* **Nginx:** Reverse proxy server
* **Cassandra:** Distributed NoSQL database

### Prerequisites

* Docker Engine (v19.03 or later) - [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
* Docker Compose (v1.29 or later) - [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)
* Golang (v1.18 or later) - [https://go.dev/doc/install](https://go.dev/doc/install)

### Running the application locally

Build and run the application with Docker Compose:

```
docker-compose up -d
```

This will build and start the Golang backend service, Redis server, Nginx reverse proxy, and Cassandra database in separate Docker containers.

```
go run main.go
```

To run GoLang file

### Contributing

We welcome contributions to this project! Please create a pull request on Github with your changes. 

**Disclaimer:** This example assumes a basic structure of the application. You may need to adjust it based on your specific implementation.
