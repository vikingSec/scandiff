# Testing Guide

This document provides instructions for running tests in the ScanDiff project.

## Backend Tests

### Running All Backend Tests

```bash
cd backend
go test ./... -v
```

### Running Specific Test Suites

#### Diff Service Tests
```bash
cd backend
go test ./services -v
```

#### API Handler Tests
```bash
cd backend
go test ./api -v
```

### Test Coverage

Generate a coverage report:

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Backend Test Structure

- **services/diff_service_test.go**: Tests for snapshot comparison logic
  - No changes detection
  - Port additions/removals
  - Software version changes
  - Vulnerability tracking
  - TLS configuration changes

- **api/snapshot_handler_test.go**: Tests for API endpoints
  - Snapshot upload
  - Invalid JSON handling
  - Missing required fields
  - Host listing
  - Health check

## Frontend Tests

### Running E2E Tests

#### Development Mode (with UI)
```bash
cd frontend
npm install
npm run test:e2e:ui
```

#### CI Mode (headless)
```bash
cd frontend
npm install
npm run test:e2e
```

### Frontend Test Structure

- **e2e/app.spec.ts**: End-to-end tests using Playwright
  - Homepage navigation
  - File upload functionality
  - Drag and drop support
  - Navigation between pages
  - Empty state handling
  - API integration

### Prerequisites for Frontend Tests

The frontend tests require:
1. Node.js 18+ installed
2. Playwright browsers installed: `npx playwright install`
3. Backend API running (or mocked)

## Integration Testing

### Manual Integration Test

1. Start the application:
```bash
docker compose up --build
```

2. Open http://localhost:3000 in your browser

3. Test the complete flow:
   - Upload a snapshot from `sample_data/`
   - View the hosts list
   - View snapshots for a host
   - Compare two snapshots
   - Verify the diff report

### Using Sample Data

The `sample_data/` directory contains real test data with three hosts:
- **125.199.235.74**: Simple IIS server with service additions
- **198.51.100.23**: Server with vulnerability changes
- **203.0.113.45**: Complex server with multiple service changes

Upload multiple snapshots for the same host to test the comparison functionality.

## Continuous Integration

For CI/CD pipelines:

```bash
# Backend tests
cd backend && go test ./... -v

# Frontend tests
cd frontend && npm install && npm run test:e2e
```

## Test Data

### Valid Snapshot Example
```json
{
  "timestamp": "2025-09-10T03:00:00Z",
  "ip": "192.168.1.1",
  "services": [
    {
      "port": 80,
      "protocol": "HTTP",
      "status": 200,
      "software": {
        "vendor": "nginx",
        "product": "nginx",
        "version": "1.22.1"
      },
      "vulnerabilities": []
    }
  ],
  "service_count": 1
}
```

### Invalid Test Cases

The tests cover various error scenarios:
- Missing required fields (IP, timestamp)
- Invalid JSON format
- Non-JSON files
- Duplicate snapshots
- Comparing snapshots from different hosts
- ID's being out of chronological order

## Debugging Tests

### Backend Debugging

Add verbose output:
```bash
go test ./... -v -run TestName
```

Run a specific test:
```bash
go test ./services -v -run TestCompareSnapshots_PortAdded
```

### Frontend Debugging

Run tests in headed mode:
```bash
cd frontend
npx playwright test --headed
```

Debug a specific test:
```bash
npx playwright test --debug -g "should display the homepage"
```

## Test Maintenance

- Keep tests focused and independent
- Use descriptive test names
- Clean up test data after tests
- Update tests when adding new features
