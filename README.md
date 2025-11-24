# 🌴 Leave Request Service (Go & Gin)

This project is a simple leave request management service built using **Go** with the **Gin** framework to provide a robust RESTful API.

## 🚀 How to Run the Project

This project utilizes a `Makefile` to simplify running various essential commands such as starting the server, performing database migrations, and running tests.

### 1\. Prerequisites

Please ensure you have the following installed in your development environment:

  * **Go** (Version $1.21+$ recommended)
  * **Docker** or **PostgreSQL** (for the database)
  * **Docker Compose** (Recommended for easy setup)
  * **`air`** (for development with hot reload - optional, but recommended)
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
  * **`migrate`** (for database migration management)
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

### 2\. Environment Configuration

1.  **Dependencies:** Fetch all Go dependencies.

    ```bash
    make tidy
    ```

2.  **Environment Variables:** Create a **`dev.env`** file in the project root and define your database and server configurations, or you can copy the `.env.sample` file into `dev.env`.

      * **For Running without Docker Compose (Host Database):**
        ```ini
        # .env
        ENV=development
        SERVER_ADDRESS=0.0.0.0:8080
        CORS_ALLOWED_ORIGIN=*
        DB_DRIVER=postgres
        DB_SOURCE=postgres://user:password@localhost:5432/leavedb?sslmode=disable //adjust
        SUPER_ADMIN_EMAIL=admin@example.com
        SUPER_ADMIN_PASSWORD=123456
        JWT_SECRET=saltandpepper
        ```
      * **For Running with Docker Compose (Recommended):**
        ```ini
        # .env
        ENV=development
        SERVER_ADDRESS=0.0.0.0:8080
        CORS_ALLOWED_ORIGIN=*
        DB_DRIVER=postgres
        DB_SOURCE=postgres://postgres:postgres@localhost:5432/leave_request?sslmode=disable
        SUPER_ADMIN_EMAIL=admin@example.com
        SUPER_ADMIN_PASSWORD=123456
        JWT_SECRET=saltandpepper
        ```
        > **Note:** **`postgres`** is the service name (hostname) for the PostgreSQL container as defined in `docker-compose.yml`.

-----

## 💻 Command Guide

### b. How to Run the Server

#### Production Mode

This command runs the server directly from the main package (`cmd/rest_api/main.go`).

```bash
make run
```

#### Development Mode

This command runs the server using `air` for automatic hot reloading upon source code changes.

```bash
make run.dev
```

### c. How to Run Migrations

#### Apply Migrations (Migrate Up)

Runs all pending migrations on the database.

```bash
make migrate
```

#### Create a New Migration File

Use the `migration` command followed by the desired migration name (replace `${name}`).

```bash
make migration name=add_initial_schema
# Example command executed:
# migrate create -ext sql -dir migration -seq add_initial_schema
```

#### Rollback the Last Migration (Migrate Down)

Reverts (rolls back) the most recent migration step.

```bash
make rollback
```

#### Seed Data

Used to populate the database with initial (dummy) data.

```bash
make seed
```

### d. How to Run Tests

#### Unit Tests

Runs all unit tests in the project. The test cache will be cleaned before execution.

```bash
make test.unit
```

#### Tests with Coverage

Runs all unit tests, generates a coverage report (`coverage.out`), and opens it in HTML format (`coverage.html`).

```bash
make test.cover
```

-----

## 🐳 Running with Docker Compose

Docker Compose allows you to launch the database service (PostgreSQL) and your Go application in an isolated environment with a single command.

### a. Start the Services (Database & API)

This command reads the `docker-compose.yml` file, builds the Go application image (if needed), and starts all defined services (e.g., `app` and `postgres`) in the background.

```bash
docker compose up -d
```

> **Application Access:** Once the services are running, your REST API application will be accessible via the configured port, typically `http://localhost:8080`.

### b. Stop the Services

This command stops and removes the running containers, networks, and volumes created by `docker compose up`.

```bash
docker compose down
```

### c. Running Migrations/Seeding against the Dockerized DB

If you want to run the `make migrate` or `make seed` commands against the database running inside the Docker container, you must ensure your local `dev.env` file is configured to point to the host's port (e.g., `localhost:5432`) if you run the commands from **outside** the container.

1.  **Stop Docker Compose:** `docker compose down` (If your services are currently running in the background).
2.  **Edit `dev.env`:** Change `DB_SOURCE` to point to `localhost`.
    ```ini
    # .env
    DATABASE_URL="postgres://postgres:postgres@localhost:5432/leave_request?sslmode=disable"
    ```
3.  **Start ONLY the Database:**
    ```bash
    docker compose up -d postgres
    ```
4.  **Run Migration/Seed (from your host machine):**
    ```bash
    make migrate
    # OR
    make seed
    ```

-----

## 📑 Brief Explanation

### a. Database Design

We have adopted a schema utilizing **PostgreSQL ENUM types** for roles and leave categories, simplifying the structure by embedding leave type and status directly into the `leave_requests` table.

#### Simple Entity-Relationship Diagram (ERD)

| Table | Key Columns | Description | PostgreSQL Type |
| :--- | :--- | :--- | :--- |
| **`users`** | `id`, `full_name`, `email`, `role` | Basic employee/user data and access level. | `role_type` ENUM |
| **`leave_requests`** | `id`, `user_id`, `start_date`, `end_date`, `type`, `status` | Details of every submitted leave request. | `leave_type_enum`, `leave_status_enum` ENUMs |

![Erd](./docs/images/ERD.png)

#### SQL Schema Details

| Type | Values | Used in Table/Column |
| :--- | :--- | :--- |
| **`role_type`** | `'superadmin'`, `'admin'`, `'employee'` | `users.role` |
| **`leave_type_enum`** | `'annual'`, `'sick'`, `'unpaid'` | `leave_requests.type` |
| **`leave_status_enum`** | `'draft'`, `'waiting_approval'`, `'approved'`, `'rejected'` | `leave_requests.status` |

#### Key Relationships

  * `users.id` $\leftrightarrow$ `leave_requests.user_id` (**One-to-Many**): A single user can have multiple leave requests.

### b. Rationale for Specific Design (Trade-off)

1.  **Using ENUM Types (e.g., `role`, `type`, `status`):**
      * **Rationale:** To enforce data integrity and restrict possible values to a predefined set directly within the database schema (e.g., a leave request *must* be one of 'annual', 'sick', or 'unpaid'). This makes querying simpler and reduces the need for join operations to look up IDs from small reference tables.
      * **Trade-off:** This design is less flexible than using separate reference tables (`leave_types`). If a new leave type needs to be added (e.g., 'maternity'), the database schema needs to be altered (a DDL operation), which is more invasive than simply inserting a new row into a `leave_types` table.

-----

## 📝 Implementation Assumptions

The following assumptions have been made to define the project's scope:

### 1\. User Roles

  * **SuperAdmin:** Possesses the highest level of access. They can perform all actions of an Admin, manage user accounts (including assigning Admin/Superadmin roles), and access system-wide configuration settings.
  * **Admin:** Responsible for approving or rejecting leave requests. They can view all leave requests across the organization.
  * **Employee Role:** Can only submit and view the status of their own leave requests.

### 2\. Leave Rules

The application logic must enforce the following rules during leave request submission and approval:

  * **Date-Time Handling and Time Zones:**
      * The system uses **date-time with time zone** (`TIMESTAMP WITH TIME ZONE`) for all date-related fields (e.g., `start_date_time`, `end_date_time`).
      * **Caution:** All date-time inputs must be stored and processed consistently, to avoid time zone conversion errors, especially when handling approvals across different geographical locations.
  * **Valid Duration:** The `start_date_time` **must not** be later than the `end_date_time`.
  * **No Past Submissions:** A leave request **cannot** be submitted if its entire duration is in the past (i.e., `end_date_time` is less than the current time).
  * **Overlap Validation (Approved Status):**
      * For the same employee, there **must not** be any two leave requests with an `APPROVED` status whose timeframes overlap.
      * *Example:* If a user has an approved leave from 2025-12-01 08:00 to 2025-12-03 17:00, no other request for that user can be approved for a period that falls within those dates/times.
  
## 🔗 API Documentation & Postman Collection

All endpoints within this service (including authentication, user management, and leave requests) can be easily tested using the provided Postman Collection.

### ⬇️ Postman Collection Link

You can import this Collection directly into Postman to begin testing:

> [[Postman Collection Link](https://dark-firefly-6451.postman.co/workspace/1beae329-aa58-46f4-8946-0d02fad55507/collection/31062578-74f67b61-d117-4bf6-b048-ab0d240f0c5b?action=share&source=copy-link&creator=31062578)]

### ⚙️ Testing Authentication Guide

* **Token Management:** The **Tests** script on the `/login` request automatically saves the received JWT token to the **`{{authToken}}`** Global/Environment variable.
* **Authorization:** All subsequent API requests that require authorization are configured to use the **Authorization Type: Bearer Token** with the variable value **`{{authToken}}`**.