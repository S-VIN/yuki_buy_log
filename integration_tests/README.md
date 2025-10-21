# Integration Tests

Comprehensive integration tests for Yuki Buy Log server using Python and pytest.

## Overview

These tests verify the complete functionality of the server including:
- User authentication and authorization
- Product CRUD operations
- Purchase CRUD operations
- Group and invite functionality
- Data sharing within groups
- Data isolation between users and groups

## Prerequisites

1. **Python 3.x** with pip
2. **PostgreSQL** database running
3. **Server** running and accessible

## Installation

Install required Python packages:

```bash
cd integration_tests
pip install -r requirements.txt
```

## Running the Server and Database

### Option 1: Using Docker Compose (Recommended for CI/CD)

```bash
# Start database and server
docker-compose up -d

# Wait for services to be ready
sleep 5

# Run tests
pytest integration_tests/
```

### Option 2: Running Manually (Recommended for Development)

#### Step 1: Start PostgreSQL Database

```bash
# Using Docker
docker run -d \
  --name yuki-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=yuki_buy_log \
  -p 5432:5432 \
  postgres:15

# Or use your local PostgreSQL installation
# Make sure the database 'yuki_buy_log' exists
```

#### Step 2: Apply Database Migrations

```bash
cd postgres
psql -h localhost -U postgres -d yuki_buy_log -f init.sql
```

#### Step 3: Start the Server

```bash
cd server

# Set database connection
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable"

# Run the server
go run .
```

The server will start on `http://localhost:8080`.

#### Step 4: Run Integration Tests

```bash
# From project root
pytest integration_tests/

# Or with verbose output
pytest integration_tests/ -v

# Or run specific test class
pytest integration_tests/test_integration.py::TestGroupSharing -v
```

## Running Tests

### Run All Tests

```bash
pytest integration_tests/
```

### Run with Verbose Output

```bash
pytest integration_tests/ -v
```

### Run Specific Test Classes

```bash
# Test authentication only
pytest integration_tests/test_integration.py::TestAuthentication -v

# Test group sharing
pytest integration_tests/test_integration.py::TestGroupSharing -v

# Test full flow
pytest integration_tests/test_integration.py::TestFullFlow -v
```

### Run with Custom Base URL

If your server is running on a different host/port:

```bash
BASE_URL=http://your-server:8080 pytest integration_tests/ -v
```

### Run with Coverage

```bash
pytest integration_tests/ --cov=integration_tests --cov-report=html
```

## Test Structure

The test suite is organized into the following test classes:

### TestServerHealth
- Server availability and basic endpoint checks

### TestAuthentication
- User registration
- User login (valid/invalid credentials)
- Access control without token
- HTTP method validation

### TestProducts
- Create products
- Get own products
- Invalid product creation

### TestPurchases
- Create purchases
- Get own purchases
- Invalid purchase creation

### TestGroupInvites
- Send invites
- Receive invites
- Mutual invite creates group
- Invite validation

### TestGroupSharing
- Product isolation before group
- Product sharing in group
- Purchase sharing in group
- Isolation for non-group users

### TestGroupManagement
- Leave group
- Group auto-deletion
- Data isolation after leaving group

### TestGroupExpansion
- Add members to existing group

### TestFullFlow
- Complete end-to-end scenarios

## Test Coverage

Current test coverage includes:

- **47 test cases** covering all major functionality
- Authentication and authorization
- Product and purchase CRUD operations
- Group formation and management
- Data sharing and isolation
- Edge cases and error handling

## Continuous Integration

These tests are designed to run in CI/CD pipelines without Docker:

```yaml
# Example GitHub Actions workflow
name: Integration Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: yuki_buy_log
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.x'

      - name: Apply database migrations
        run: |
          psql -h localhost -U postgres -d yuki_buy_log -f postgres/init.sql
        env:
          PGPASSWORD: postgres

      - name: Start server
        run: |
          cd server
          go run . &
          sleep 5
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable

      - name: Install test dependencies
        run: pip install -r integration_tests/requirements.txt

      - name: Run integration tests
        run: pytest integration_tests/ -v
```

## Troubleshooting

### Connection Refused Error

If you get "Connection refused" errors:

1. Make sure the server is running: `curl http://localhost:8080/products`
2. Check the server logs for errors
3. Verify the database is running and accessible

### Database Errors

If tests fail with database errors:

1. Check PostgreSQL is running: `pg_isready -h localhost`
2. Verify database exists: `psql -h localhost -U postgres -l | grep yuki_buy_log`
3. Check migrations are applied: Look for `users`, `products`, `purchases`, `groups`, `invites` tables

### Test Isolation Issues

Each test creates unique users with UUIDs to avoid conflicts. If you see test failures:

1. Check if the database needs to be reset
2. Restart the server to clear any in-memory state
3. Run tests individually to identify the failing test

### Server Not Starting

If the server fails to start:

1. Check if port 8080 is already in use: `lsof -i :8080`
2. Verify DATABASE_URL is set correctly
3. Check Go dependencies are installed: `cd server && go mod download`

## Development Workflow

For local development:

```bash
# Terminal 1: Start database
docker run --rm -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=yuki_buy_log \
  postgres:15

# Terminal 2: Start server
cd server
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable"
go run .

# Terminal 3: Run tests
pytest integration_tests/ -v

# Or run in watch mode (requires pytest-watch)
ptw integration_tests/
```

## Contributing

When adding new features:

1. Write integration tests first
2. Run tests locally to verify
3. Ensure all tests pass before committing
4. Update this README if adding new test categories

## Notes

- Tests use unique user IDs (UUIDs) to avoid conflicts
- Each test is independent and can run in parallel
- Tests assume a clean server state (database can have existing data)
- Some tests create groups and verify sharing - these cannot be parallelized
