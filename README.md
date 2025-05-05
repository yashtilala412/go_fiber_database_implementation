echo '# Golang BoilerPlater

## Index
- [Introduction](#introduction)
- [Add Storage Drivers](#add-storage-drivers)
- [Multiple Database System Support](#multiple-database-system-support)
- [Kick Start Commands](#kick-start-commands)
- [Migrations](#migrations)
- [Database Seeding](#database-seeding)
- [Kratos Integration](#kratos-integration)
- [Messaging Queue](#messaging-queue)
- [Code Walk-through](#code-walk-through)
  - [Config](#config)
  - [Command](#command)
  - [Route](#route)
  - [Middleware](#middleware)
  - [Model](#model)
  - [Controller](#controller)
  - [Utils](#utils)
- [Testcases](#testcases)

---

## Introduction

- This template is built on [Fiber](https://github.com/gofiber/fiber).
- Copy `.env.example` to `.env` for environment configuration.
- To sync Go packages:

  \`\`\`bash
  go mod vendor
  \`\`\`

- To add a new package:

  \`\`\`bash
  go get {package_url}
  go mod vendor
  \`\`\`

- To remove unused packages:

  \`\`\`bash
  go mod tidy
  \`\`\`

- Make sure you set proper DB configuration in your `.env` file.

---

## Add Storage Drivers

- Reference: [Fiber Storage Documentation](https://docs.gofiber.io/storage/)

---

## Multiple Database System Support

- Supported Databases: `postgres`, `mysql`, and `sqlite3`
- Allows switching DB systems with minimal code changes (some adjustments may be needed).
- Uses [`goqu`](https://github.com/doug-martin/goqu) SQL builder.
- Set the DB system using the `DB_DIALECT` variable:

  \`\`\`env
  DB_DIALECT=postgres
  \`\`\`

---

## Kick Start Commands

These commands are defined in the `Makefile`.

> 🔧 Make sure to update the `Makefile` if you change folder structure (e.g. `app.go` location).

- `make start-api-dev` – Starts app with `nodemon` for live-reload.
- `make start-api` – Runs app using `go run app.go`.
- `make migrate file_name={MIGRATION_FILE_NAME}` – Creates both `up` and `down` migrations.
- `make build app_name={BINARY_NAME}` – Builds a binary of your project.
- `make install app_name={BINARY_NAME}` – Builds an optimized binary using `-s -w` flags.
- `make test` – Runs all test cases (uses `.env.testing`).
- `make test-wo-cache` – Runs test cases without caching.
- `make swagger-gen` – Generates Swagger docs.
- `make migrate-up` – Runs all `up` migrations.

---

## Migrations

This project supports auto-generation and execution of migrations.

Create a new migration:

\`\`\`bash
make migrate file_name=create_users_table
\`\`\`

Run all `up` migrations:

\`\`\`bash
make migrate-up
\`\`\`

---

## Database Seeding

This project supports database seeding with initial data, which can be loaded from CSV or other static data files. Follow the steps below to seed your database:

### Steps for Database Seeding

1. **Prepare your seed files**:
    - Place your CSV files or static data in the appropriate directory (e.g., `./seeds/`).

2. **Run the seeding command**:

    The following command will execute the seeding process, loading data into your database:

    ```bash
    go run app.go seed
    ```

    This command will look for the seed files in the appropriate directory and load the data accordingly.

    > **Note:** Ensure that the directory structure is correct and the seed files are properly formatted for the database schema.

3. **Download dataset (if needed)**:
    If your project requires a specific dataset to seed the database, you can download it by running a command like:

    ```bas[h
   (https://www.kaggle.com/datasets/lava18/google-play-store-apps)
    ```

    After downloading the dataset, run the seeding command:

    ```bash
    go run app.go seed
    ```

> This will populate your database with the initial dataset from the CSV or static files.

---

## Kratos Integration

- This project is integrated with [ORY Kratos](https://www.ory.sh/kratos/) for user authentication.
- You can configure endpoints and middleware to sync with Kratos flows.

---

## Messaging Queue

- The template includes integration with messaging queues like NATS or RabbitMQ.
- Configure in `.env` and initialize in `pkg/queue`.

---

## Code Walk-through

### Config

- All environment variables and app settings are managed in `pkg/config`.

### Command

- CLI commands and task runners are in `cmd/`.

### Route

- Routes are defined in `pkg/routes`.

### Middleware

- Custom middlewares are located in `pkg/middleware`.

### Model

- Database models live in `pkg/models`.

### Controller

- Business logic and handlers are in `pkg/controller`.

### Utils

- Helper functions and utilities are in `pkg/utils`.

---

## Testcases

- Test files use Go’s built-in `testing` package.
- Use the `.env.testing` file for test-specific configuration.
- To run all tests:

  \`\`\`bash
  make test
  \`\`\`

- Without cache:

  \`\`\`bash
  make test-wo-cache
  \`\`\`
' > README.md
