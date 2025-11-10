# Scripts Documentation

This directory contains automation scripts for managing the E-Wallet application.

## Quick Start

```bash
# Start full development environment
./scripts/manage.sh dev

# Run API tests
./scripts/manage.sh test

# Generate HMAC signature
./scripts/manage.sh hmac
```

---

## manage.sh

Main management script for the E-Wallet application. Handles Docker, database, application lifecycle, and testing.

### Docker Commands

#### `./scripts/manage.sh up`
Starts Docker containers (PostgreSQL and Redis).

**What it does:**
- Starts `docker-compose up -d`
- Waits for PostgreSQL (5 seconds)
- Waits for Redis (2 seconds)
- Shows container status

**Example:**
```bash
./scripts/manage.sh up
```

#### `./scripts/manage.sh down`
Stops all Docker containers.

**What it does:**
- Runs `docker-compose down`
- Stops PostgreSQL and Redis containers

**Example:**
```bash
./scripts/manage.sh down
```

#### `./scripts/manage.sh restart`
Restarts all Docker containers.

**What it does:**
- Runs `docker-compose restart`
- Shows container status after restart

**Example:**
```bash
./scripts/manage.sh restart
```

#### `./scripts/manage.sh logs [service]`
Views Docker container logs.

**Parameters:**
- `service` (optional) - Service name (postgres, redis)

**What it does:**
- Shows last 100 lines of logs
- Follows log output (live updates)
- If no service specified, shows all logs

**Examples:**
```bash
./scripts/manage.sh logs              # All services
./scripts/manage.sh logs postgres     # PostgreSQL only
./scripts/manage.sh logs redis        # Redis only
```

#### `./scripts/manage.sh status`
Shows status of all Docker containers.

**What it does:**
- Runs `docker-compose ps`
- Shows container names, status, ports

**Example:**
```bash
./scripts/manage.sh status
```

---

### Application Commands

#### `./scripts/manage.sh dev`
**Development mode - local app with Docker DBs (RECOMMENDED)**

**What it does:**
1. Starts Docker containers (PostgreSQL + Redis only)
2. Waits for PostgreSQL to be ready (5 seconds)
3. Initializes database (creates `e_wallet_db` + extensions)
4. Truncates existing data
5. Seeds database with test data
6. Starts application locally with `go run cmd/server/main.go`

**Example:**
```bash
./scripts/manage.sh dev
```

**Output:**
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Health: http://localhost:8080/health

#### `./scripts/manage.sh prod`
**Production mode - all services in Docker**

**What it does:**
1. Builds Docker image for application
2. Starts all services (PostgreSQL + Redis + App)
3. Waits for health check (10 seconds)
4. Verifies application is running

**Example:**
```bash
./scripts/manage.sh prod
```

**Output:**
- API: http://localhost:8080
- All services running in Docker
- Auto-restart enabled

#### `./scripts/manage.sh build`
Builds application binary.

**What it does:**
- Compiles Go code
- Creates binary at `bin/e-wallet`

**Example:**
```bash
./scripts/manage.sh build
```

#### `./scripts/manage.sh run`
Runs application from binary.

**What it does:**
- Checks if `bin/e-wallet` exists
- If not, builds it first
- Runs the binary

**Example:**
```bash
./scripts/manage.sh run
```

#### `./scripts/manage.sh test`
Runs API test suite.

**What it does:**
- Checks if application is running (http://localhost:8080/health)
- Executes `scripts/test-api.sh`
- Tests all 4 API endpoints + authentication

**Example:**
```bash
# Start app first
./scripts/manage.sh dev

# In another terminal
./scripts/manage.sh test
```

**Tests executed:**
1. Check wallet existence (valid)
2. Check non-existent wallet
3. Get wallet balance
4. Deposit to wallet
5. Get monthly statistics
6. Deposit exceeding limit (should fail)
7. Invalid amount (should fail)
8. Missing authentication (should fail)
9. Invalid HMAC signature (should fail)

#### `./scripts/manage.sh hmac`
Runs HMAC signature generator tool.

**What it does:**
- Launches interactive HMAC generator
- Helps create authenticated API requests

**Example:**
```bash
./scripts/manage.sh hmac
```

---

### Database Commands

#### `./scripts/manage.sh init`
Initializes database.

**What it does:**
- Creates `e_wallet_db` database if not exists
- Connects to the database
- Tables are auto-created by GORM on app start

**Example:**
```bash
./scripts/manage.sh init
```

**SQL executed:**
```sql
SELECT 'CREATE DATABASE e_wallet_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'e_wallet_db')\gexec

\c e_wallet_db
```

#### `./scripts/manage.sh seed`
Seeds database with test data.

**What it does:**
- Truncates existing data (api_clients, wallets, transactions)
- Inserts test API clients
- Inserts test wallets (identified + unidentified)
- Resets auto-increment sequences

**Example:**
```bash
./scripts/manage.sh seed
```

**Test data inserted:**
- **API Clients:**
  - `alif_partner` / `alif_secret_2025`
  - `megafon_api` / `megafon_key_secure`
  - `tcell_integration` / `tcell_hmac_key`

- **Unidentified Wallets (max 10,000 TJS):**
  - `992900123456` - 2,500 TJS
  - `992935789012` - 5,000 TJS
  - `992918765432` - 7,500 TJS
  - `992987654321` - 10,000 TJS
  - `992901234567` - 0 TJS

- **Identified Wallets (max 100,000 TJS):**
  - `992900111222` - 25,000 TJS
  - `992935333444` - 50,000 TJS
  - `992918555666` - 75,000 TJS
  - `992987777888` - 100,000 TJS
  - `992901112233` - 0 TJS

#### `./scripts/manage.sh backup`
Creates database backup.

**What it does:**
- Creates `backups/` directory if not exists
- Dumps database to `backup_YYYYMMDD_HHMMSS.sql`
- Compresses with gzip
- Keeps only last 7 backups

**Example:**
```bash
./scripts/manage.sh backup
```

**Output:**
```
backups/backup_20251110_141500.sql.gz
```

#### `./scripts/manage.sh restore <backup_file>`
Restores database from backup.

**Parameters:**
- `backup_file` - Path to backup file (.sql or .sql.gz)

**What it does:**
- Decompresses if .gz
- Restores to `e_wallet_db`

**Example:**
```bash
./scripts/manage.sh restore backups/backup_20251110_141500.sql.gz
```

---

### Utility Commands

#### `./scripts/manage.sh swagger`
Generates Swagger documentation.

**What it does:**
- Checks if `swag` is installed
- Parses Go annotations
- Generates `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`

**Example:**
```bash
./scripts/manage.sh swagger
```

**Requirements:**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**View documentation:**
http://localhost:8080/swagger/index.html (only in development mode)

#### `./scripts/manage.sh clean`
Cleans up everything.

**What it does:**
- Stops and removes Docker containers
- Removes Docker volumes
- Removes `bin/e-wallet` binary
- Removes logs from `logs/` directory

**Example:**
```bash
./scripts/manage.sh clean
```

**Warning:** This will delete all data in Docker volumes!

---

## test-api.sh

Automated API test suite. Tests all endpoints with HMAC authentication.

### Usage

```bash
./scripts/test-api.sh
```

### Configuration

Edit these variables at the top of the script:

```bash
API_URL="http://localhost:8080/api/v1"
USER_ID="alif_partner"
SECRET_KEY="alif_secret_2025"
```

### Test Cases

1. **Check wallet existence** - Valid wallet
2. **Check non-existent wallet** - Should return exists=false
3. **Get wallet balance** - Returns balance in dirams
4. **Deposit to wallet** - Deposits 100 TJS (10,000 dirams)
5. **Get monthly statistics** - Returns count and total amount
6. **Deposit exceeding limit** - Should fail with error
7. **Invalid amount** - Negative amount should fail
8. **Missing authentication** - No headers should fail
9. **Invalid HMAC signature** - Wrong digest should fail

### Output

```
=========================================
E-Wallet API Test Suite
=========================================

Test 1: Check Wallet Existence
Request: {"account_id":"992900123456"}
Response: {"exists":true,"account_id":"992900123456"}
✓ Test passed
=========================================
...
```

---

## Environment Variables

The `manage.sh` script uses these environment variables (with defaults):

```bash
DB_HOST=localhost              # PostgreSQL host
DB_PORT=5432                   # PostgreSQL port
DB_USER=postgres               # PostgreSQL user
DB_NAME=e_wallet_db            # Database name
POSTGRES_PASSWORD=postgres     # PostgreSQL password
```

These match the `.env` file variables:

```bash
# .env
APP_ENVIRONMENT=development
POSTGRES_PASSWORD=postgres
REDIS_PASSWORD=
```

### Override example:

```bash
POSTGRES_PASSWORD=mypass ./scripts/manage.sh seed
```

---

## Common Workflows

### First Time Setup

```bash
# 1. Start everything
./scripts/manage.sh dev

# 2. In another terminal, test API
./scripts/manage.sh test

# 3. Generate HMAC for manual testing
./scripts/manage.sh hmac
```

### Daily Development

```bash
# Start dev environment
./scripts/manage.sh dev

# View logs if needed
./scripts/manage.sh logs

# Stop when done
./scripts/manage.sh down
```

### Manual Setup (Step by Step)

```bash
# 1. Start Docker
./scripts/manage.sh up

# 2. Initialize database
./scripts/manage.sh init

# 3. Seed test data
./scripts/manage.sh seed

# 4. Build and run
./scripts/manage.sh build
./scripts/manage.sh run
```

### Testing

```bash
# Run automated tests
./scripts/manage.sh test

# Or test manually with HMAC generator
./scripts/manage.sh hmac
```

### Database Management

```bash
# Backup database
./scripts/manage.sh backup

# Restore from backup
./scripts/manage.sh restore backups/backup_20251110_141500.sql.gz

# Re-seed with fresh data
./scripts/manage.sh seed
```

### Cleanup

```bash
# Stop containers
./scripts/manage.sh down

# Full cleanup (removes volumes and data)
./scripts/manage.sh clean
```

---

## Troubleshooting

### Port already in use

```bash
# Check what's using port 8080
lsof -ti:8080

# Kill the process
lsof -ti:8080 | xargs kill -9

# Or use different port in config
```

### PostgreSQL connection refused

```bash
# Check if PostgreSQL is running
./scripts/manage.sh status

# View PostgreSQL logs
./scripts/manage.sh logs postgres

# Restart PostgreSQL
./scripts/manage.sh restart
```

### Database doesn't exist

```bash
# Initialize database
./scripts/manage.sh init

# Or use dev command (does everything)
./scripts/manage.sh dev
```

### Tests failing

```bash
# Make sure app is running
curl http://localhost:8080/health

# Check if database is seeded
./scripts/manage.sh seed

# View app logs
./scripts/manage.sh logs
```

---

## Script Permissions

Make scripts executable:

```bash
chmod +x scripts/manage.sh
chmod +x scripts/test-api.sh
```

---

## Dependencies

### Required

- Docker & Docker Compose
- Go 1.21+
- PostgreSQL client (`psql`)
- OpenSSL (for HMAC generation in tests)

### Optional

- `swag` (for Swagger generation)
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

---

## Notes

- All amounts are in **dirams** (1 TJS = 100 dirams)
- Unidentified wallets: max 10,000 TJS (1,000,000 dirams)
- Identified wallets: max 100,000 TJS (10,000,000 dirams)
- HMAC-SHA1 authentication required for all API endpoints
- Swagger documentation only available in development mode
- Redis caching enabled for API client authentication
