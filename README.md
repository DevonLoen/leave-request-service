# 🌴 Leave Request Service (Go & Gin)

This project is a simple leave request management service built using **Go** with the **Gin** framework to provide a robust RESTful API.

## 🚀 How to Run the Project

This project utilizes a `Makefile` to simplify running various essential commands such as starting the server, performing database migrations, and running tests.

### 1\. Prerequisites

Please ensure you have the following installed in your development environment:

  * **Go** (Version $1.21+$ recommended)
  * **Docker** or **PostgreSQL** (for the database)
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
2.  **Environment Variables:** Create a `dev.env` file in the project root and define your database and server configurations or you can copy .env.sample file into dev.env file, for example:
    ```ini
    # .env
    DATABASE_URL="postgres://user:password@localhost:5432/leavedb?sslmode=disable"
    SERVER_PORT=8080
    ```

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
make run dev
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

## 📑 Brief Explanation

### a. Database Design

We have adopted a schema utilizing **PostgreSQL ENUM types** for roles and leave categories, simplifying the structure by embedding leave type and status directly into the `leave_requests` table.

#### Simple Entity-Relationship Diagram (ERD)

| Table | Key Columns | Description | PostgreSQL Type |
| :--- | :--- | :--- | :--- |
| **`users`** | `id`, `full_name`, `email`, `role` | Basic employee/user data and access level. | `role_type` ENUM |
| **`leave_requests`** | `id`, `user_id`, `start_date`, `end_date`, `type`, `status` | Details of every submitted leave request. | `leave_type_enum`, `leave_status_enum` ENUMs |

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

