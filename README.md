## go-testcontainers

### Table Of Contents

- [What are Testcontainers?](#what-are-testcontainers)
- [Why use Testcontainers over regular docker and docker compose  ?](#why-use-testcontainers-over-regular-docker-and-docker-compose-)
- [Benefits of Using Testcontainers ?](#benefits-of-using-testcontainers-)
- [Running the Project](#running-the-project)

### What are Testcontainers?

Testcontainers is a library that provides easy and lightweight APIs for bootstrapping local development and test dependencies with real services wrapped in Docker containers. Using Testcontainers, you can write tests that depend on the same services you use in production without mocks or in-memory services.

### Why use Testcontainers over regular docker and docker compose  ?

Testcontainers offer a more automated and consistent approach to managing Docker containers for testing compared to regular development containers. They ensure that each test runs in a clean, isolated environment by automatically starting and stopping containers as needed, which eliminates the risk of environment drift or contamination between test runs.

### Benefits of Using Testcontainers ?

1. On-demand infrastructure: Testcontainers create isolated environments for each test, avoiding pre-provisioning and preventing data issues in parallel tests.


2. Consistent testing: Run integration tests easily from your IDE, with the same experience locally and in CI.


3. Reliable setup: Testcontainers ensure services are ready before tests start, using built-in wait strategies.


4. Networking support: Testcontainers manage port mappings and allow containers to communicate within a Docker network.


5. Automatic cleanup: Resources are automatically cleaned up after tests, even if the test crashes.

### Running the Project

1. Clone the project and navigate to the project directory. 

```
git clone git@github.com:pgaijin66/go-testcontainers.git
cd go-testcontainers
```

2. Install dependencies

```
go mod tidy
```

3. Execute unit tests using Testcontainters

```
$ task test

task: [build] go build -o app main.go
task: [test] go test -v ./...
=== RUN   TestGetUsersHandler
=== RUN   TestGetUsersHandler/Handle_multiple_users
2024/09/02 08:41:04 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 26.1.1
  API Version: 1.45
  Operating System: Docker Desktop
  Total Memory: 7840 MB
  Labels:
    com.docker.desktop.address=unix:///Users/pthapa/Library/Containers/com.docker.docker/Data/docker-cli.sock
  Testcontainers for Go Version: v0.33.0
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: b68b3045cc70383fc283797770c787a3f7cb5ca078ce28e9b95ce78b290caad4
  Test ProcessID: 0605e174-1000-478d-b17e-cacc9be55c2a
2024/09/02 08:41:04 🐳 Creating container for image testcontainers/ryuk:0.8.1
2024/09/02 08:41:04 ✅ Container created: 5a4326a9432c
2024/09/02 08:41:04 🐳 Starting container: 5a4326a9432c
2024/09/02 08:41:04 ✅ Container started: 5a4326a9432c
2024/09/02 08:41:04 ⏳ Waiting for container id 5a4326a9432c image: testcontainers/ryuk:0.8.1. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false}
2024/09/02 08:41:04 🔔 Container is ready: 5a4326a9432c
2024/09/02 08:41:04 🐳 Creating container for image postgres:13
2024/09/02 08:41:04 ✅ Container created: 3a3547b2adb4
2024/09/02 08:41:04 🐳 Starting container: 3a3547b2adb4
2024/09/02 08:41:05 ✅ Container started: 3a3547b2adb4
2024/09/02 08:41:05 ⏳ Waiting for container id 3a3547b2adb4 image: postgres:13. Waiting for: &{timeout:<nil> deadline:<nil> Strategies:[0x1400046fbc0 0x1400046fbf0]}
2024/09/02 08:41:05 🔔 Container is ready: 3a3547b2adb4
2024/09/02 08:41:05 🐳 Terminating container: 3a3547b2adb4
2024/09/02 08:41:05 🚫 Container terminated: 3a3547b2adb4
--- PASS: TestGetUsersHandler (1.34s)
    --- PASS: TestGetUsersHandler/Handle_multiple_users (1.34s)
PASS
ok   my-app 1.965s
```