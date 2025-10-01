# Getting Started with ScanDiff

## Prerequisites

Before you begin, make sure you have:
- **Docker** (version 20.10+) - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** (version 2.0+) - Usually included with Docker Desktop

## Quick Start (3 Steps)

### Step 1: Start the Application

Open your terminal and run:

```bash
./quickstart.sh
```

Or manually:

```bash
docker compose up --build
```

This will:
- Build the backend (Go) container
- Build the frontend (React) container
- Start both services
- Initialize the SQLite database
- Set up networking between containers

Wait for the message: "🎉 ScanDiff is now running!"

The result of this command will be two networked containers. The API can be reached via `http://localhost:8080/api/*` and the frontend UI can be reached via `http://localhost:3000`

### Step 2: Upload Your First Snapshot

1. Click the upload area or drag a JSON file
2. Use one of the sample files from `sample_data/` folder
3. For example: `sample_data/host_125.199.235.74_2025-09-10T03-00-00Z.json`
4. Click "Upload Snapshot"
5. You'll see a success message!

## Next Steps

### Upload More Snapshots

To see the comparison feature in action, upload at least 2 snapshots for the same host:

1. Upload `host_125.199.235.74_2025-09-10T03-00-00Z.json`
2. Upload `host_125.199.235.74_2025-09-15T08-49-45Z.json`
3. Click "Hosts" in the navigation
4. Click on the host IP "125.199.235.74"
5. Select both snapshots (checkboxes)
6. Click "Compare Snapshots"
7. View the detailed diff report!


## Common Tasks

### Stop the Application
```bash
docker compose down
```

### View Logs
```bash
docker compose logs -f
```

### Restart the Application
```bash
docker compose restart
```

### Clean Everything (including database)
```bash
docker compose down -v
rm -rf data/
```

### Run Tests
```bash
# Backend tests
cd backend && go test ./... -v

# Frontend tests (requires npm install first)
cd frontend && npm install && npm run test:e2e
```

## Troubleshooting

### Port Already in Use

If you see errors about ports 3000 or 8080 being in use:

**Option 1**: Stop other services using these ports

**Option 2**: Change ports in `docker-compose.yml`:
```yaml
services:
  frontend:
    ports:
      - "3001:3000"  # Changed from 3000
  backend:
    ports:
      - "8081:8080"  # Changed from 8080
```

### Docker Not Running

Make sure Docker Desktop is running:
- **Mac**: Check the Docker icon in the menu bar
- **Windows**: Check the Docker icon in the system tray
- **Linux**: Run `sudo systemctl start docker`

### Application Won't Start

1. Check Docker logs:
   ```bash
   docker compose logs
   ```

2. Try rebuilding:
   ```bash
   docker compose down
   docker compose up --build
   ```

3. Check if containers are running:
   ```bash
   docker compose ps
   ```

### Upload Fails

Make sure you're uploading valid JSON files with required fields:
- `ip` (string)
- `timestamp` (ISO-8601 format)
- `services` (array)
- `service_count` (number)

See the sample data files for examples.

### Can't Compare Snapshots

You can only compare snapshots from the **same host**. Make sure:
- Both snapshots have the same IP address
- You've selected exactly 2 snapshots
- The snapshots are different (different timestamps)

## API Usage (Optional)

If you want to use the API directly:

### Upload a Snapshot
```bash
curl -X POST http://localhost:8080/api/snapshots/upload \
  -F "file=@sample_data/host_125.199.235.74_2025-09-10T03-00-00Z.json"
```

### List All Hosts
```bash
curl http://localhost:8080/api/hosts
```

### List Snapshots for a Host
```bash
curl http://localhost:8080/api/hosts/125.199.235.74/snapshots
```

### Compare Two Snapshots
```bash
curl http://localhost:8080/api/diff/1/2
```

### Health Check
```bash
curl http://localhost:8080/health
```

