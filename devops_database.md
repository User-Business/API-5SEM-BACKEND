# 📌 Database DevOps Implementation

This PR implements the Database DevOps structure for the project using Docker, PostgreSQL and Atlas Schema Migration.

The main goal was to standardize the database environment, automate database initialization and introduce schema versioning with migration integrity validation.

---

# 🚀 Technologies Used

* Docker
* Docker Compose
* PostgreSQL 17
* Atlas Schema Migration
* pgAdmin

---

# ⚙️ Main Implementations

## 🐳 Database Environment with Docker

The PostgreSQL container was configured using Docker Compose to ensure a standardized environment for all developers.

### Configured Services

* PostgreSQL
* Backend
* Frontend
* ETL
* pgAdmin

### Main Command Used

```bash
docker compose up -d
```

This command initializes the entire development environment automatically.

---

# 🗄️ PostgreSQL Configuration

Environment variables were configured using the `.env` file.

```env
POSTGRES_USER=denarius
POSTGRES_PASSWORD=denarius
POSTGRES_DB=denarius
```

Database connection string:

```env
DATABASE_URL=postgres://denarius:denarius@db:5432/denarius?sslmode=disable
```

---

# 📂 Database Migrations

SQL migrations were organized inside the following directory:

```text
migrations/
```

Main migration file:

```text
001_create_tables.sql
```

Responsible for creating:

* transactional tables
* Data Warehouse tables
* constraints
* foreign keys
* schemas
* database structure

---

# 🧩 Atlas Integration

Atlas was installed and configured locally for schema inspection and migration management.

## Atlas Configuration

Created file:

```text
atlas.hcl
```

Configuration:

```hcl
env "local" {
  url = "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable"

  dev = "docker://postgres/17/dev"

  migration {
    dir = "file://migrations"
  }
}
```

---

# 📌 Schema Versioning

Database schema inspection was performed using:

```bash
C:\atlas\atlas.exe schema inspect -u "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable" > schema.hcl
```

Generated file:

```text
schema.hcl
```

This file stores:

* tables
* columns
* constraints
* relationships
* schemas

The `schema.hcl` file represents the current database structure and allows schema version tracking.

---

# 🔒 Migration Integrity Validation

Migration integrity validation was configured using:

```bash
C:\atlas\atlas.exe migrate hash --dir "file://migrations"
```

Generated file:

```text
atlas.sum
```

Used for:

* migration integrity validation
* migration tracking
* schema consistency verification

This ensures that migrations are not accidentally modified after being versioned.

---

# 👥 Team Database Standardization

One of the main objectives of this implementation is ensuring that every developer uses the exact same database structure during development.

## Development Flow

When a developer clones the repository:

### 1. Start Containers

```bash
docker compose up -d
```

### 2. PostgreSQL Container Starts Automatically

The database container is initialized using the project configuration.

### 3. SQL Migrations Are Executed

The migration scripts inside:

```text
migrations/
```

are automatically applied.

### 4. Database Structure Is Standardized

All developers work using:

* same PostgreSQL version
* same tables
* same constraints
* same schema structure

This eliminates environment inconsistencies across the team.

---

# 🔄 Future Database Changes

When new database changes are required:

### 1. Create a New Migration

Example:

```text
002_add_new_column.sql
```

### 2. Inspect the Updated Schema

```bash
C:\atlas\atlas.exe schema inspect -u "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable" > schema.hcl
```

### 3. Update Migration Hash

```bash
C:\atlas\atlas.exe migrate hash --dir "file://migrations"
```

### 4. Commit Files

Developers must commit:

* migration SQL files
* `schema.hcl`
* `atlas.sum`

This guarantees version synchronization across the entire team.

---

# 📊 Database DevOps Flow

```text
Developer
    ↓
SQL Migration
    ↓
GitHub Repository
    ↓
Docker Compose
    ↓
PostgreSQL Container
    ↓
Atlas Schema Inspection
    ↓
schema.hcl
    ↓
atlas.sum
    ↓
Database Versioning
    ↓
Team Synchronization
```

---

# ✅ Benefits

* Standardized database environment
* Automated database initialization
* Schema versioning
* Migration integrity validation
* Easier onboarding for developers
* Better DevOps workflow
* Improved team synchronization
* Future-ready migration structure

---

# 📁 Final Structure

```text
API-5SEM-BACKEND/
│
├── atlas.hcl
├── schema.hcl
│
├── migrations/
│   ├── 001_create_tables.sql
│   └── atlas.sum
```

---

# 🧪 Validation

Commands successfully executed:

```bash
docker compose up -d

docker ps

C:\atlas\atlas.exe version

C:\atlas\atlas.exe schema inspect -u "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable" > schema.hcl

C:\atlas\atlas.exe migrate hash --dir "file://migrations"
```