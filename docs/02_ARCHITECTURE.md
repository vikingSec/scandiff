# Architecture Documentation

## System Overview

ScanDiff is a full-stack web application for tracking and analyzing security changes in network hosts over time. The system follows a three-tier architecture:

1. **Presentation Layer**: React-based web interface
2. **Application Layer**: Go/Gin REST API
3. **Data Layer**: SQLite database

## Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend (React)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Upload  │  │  Hosts   │  │Snapshots │  │   Diff   │   │
│  │   Page   │  │   Page   │  │   Page   │  │   Page   │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│         │              │              │              │       │
│         └──────────────┴──────────────┴──────────────┘       │
│                          │                                    │
│                    ┌─────▼─────┐                             │
│                    │  API Client│                             │
│                    └─────┬─────┘                             │
└──────────────────────────┼───────────────────────────────────┘
                           │ HTTP/JSON
                           │
┌──────────────────────────▼───────────────────────────────────┐
│                    Backend (Go/Gin)                          │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    API Handlers                      │    │
│  │  ┌──────────────────┐  ┌──────────────────┐        │    │
│  │  │ Snapshot Handler │  │   Diff Handler   │        │    │
│  │  └────────┬─────────┘  └────────┬─────────┘        │    │
│  └───────────┼─────────────────────┼──────────────────┘    │
│              │                     │                         │
│  ┌───────────▼─────────────────────▼──────────────────┐    │
│  │                    Services                          │    │
│  │  ┌──────────────────┐  ┌──────────────────┐        │    │
│  │  │ Snapshot Service │  │   Diff Service   │        │    │
│  │  └────────┬─────────┘  └──────────────────┘        │    │
│  └───────────┼──────────────────────────────────────────┘    │
│              │                                                │
│  ┌───────────▼──────────────────────────────────────────┐    │
│  │                  Database Layer                       │    │
│  │              (SQLite Connection)                      │    │
│  └───────────┬──────────────────────────────────────────┘    │
└──────────────┼───────────────────────────────────────────────┘
               │
┌──────────────▼───────────────────────────────────────────────┐
│                       SQLite Database                         │
│  ┌───────────────────────────────────────────────────────┐   │
│  │                    snapshots Table                     │   │
│  │  - id (PK)                                            │   │
│  │  - ip                                                  │   │
│  │  - timestamp                                           │   │
│  │  - filename                                            │   │
│  │  - data (JSON)                                        │   │
│  │  - uploaded_at                                        │   │
│  └───────────────────────────────────────────────────────┘   │
└───────────────────────────────────────────────────────────────┘
```

## Backend Architecture

### Directory Structure

```
backend/
├── main.go                 # Application entry point
├── api/                    # HTTP handlers and routing
│   ├── router.go          # Route definitions
│   ├── snapshot_handler.go # Snapshot endpoints
│   ├── diff_handler.go    # Diff endpoints
│   └── *_test.go          # Handler tests
├── services/              # Business logic
│   ├── snapshot_service.go # Snapshot CRUD operations
│   ├── diff_service.go    # Comparison logic
│   └── *_test.go          # Service tests
├── models/                # Data structures
│   ├── snapshot.go        # Snapshot models
│   └── diff.go            # Diff models
├── db/                    # Database layer
│   └── db.go              # Database initialization
├── Dockerfile             # Container definition
├── go.mod                 # Dependencies
└── go.sum                 # Dependency checksums
```

### Key Components

#### 1. API Layer (`api/`)

**Router (`router.go`)**
- Sets up HTTP routes
- Configures CORS
- Initializes middleware
- Connects handlers to services

**Snapshot Handler (`snapshot_handler.go`)**
- `POST /api/snapshots/upload`: Upload snapshot file
- `GET /api/snapshots/:id`: Get snapshot by ID
- `GET /api/hosts`: List all hosts
- `GET /api/hosts/:ip/snapshots`: List snapshots for a host

**Diff Handler (`diff_handler.go`)**
- `GET /api/diff/:id1/:id2`: Compare two snapshots

#### 2. Service Layer (`services/`)

**Snapshot Service**
- `CreateSnapshot()`: Store new snapshot
- `GetSnapshotByID()`: Retrieve snapshot
- `ListSnapshotsByIP()`: Get all snapshots for a host
- `ListAllHosts()`: Get all unique IPs

**Diff Service**
- `CompareSnapshots()`: Generate diff report
- `compareServices()`: Compare individual services
- `compareTLS()`: Compare TLS configurations
- `compareVulnerabilities()`: Track CVE changes

#### 3. Model Layer (`models/`)

**Data Structures**
- `Snapshot`: Host snapshot with services
- `Service`: Network service details
- `Software`: Software information
- `TLS`: TLS configuration
- `DiffReport`: Comparison result
- `PortChange`: Port addition/removal
- `ServiceChange`: Service modification

#### 4. Database Layer (`db/`)

**Schema**
```sql
CREATE TABLE snapshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    filename TEXT NOT NULL,
    data TEXT NOT NULL,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ip, timestamp)
);

CREATE INDEX idx_snapshots_ip ON snapshots(ip);
CREATE INDEX idx_snapshots_timestamp ON snapshots(timestamp);
```

## Frontend Architecture

### Directory Structure

```
frontend/
├── src/
│   ├── main.tsx           # Application entry
│   ├── App.tsx            # Root component & routing
│   ├── index.css          # Global styles (Tailwind)
│   ├── pages/             # Page components
│   │   ├── UploadPage.tsx # Upload snapshots
│   │   ├── HostsPage.tsx  # View all hosts
│   │   ├── SnapshotsPage.tsx       # View all snapshots for a host
|   |   ├── SnapshotDetailPage.tsx  # View details on a snapshot
│   │   └── DiffPage.tsx            # View information on a diff between two snapshots
│   ├── api/               # API client
│   │   └── api.ts
│   └── types/             # TypeScript types
│       └── index.ts
├── e2e/                   # Playwright tests
│   └── app.spec.ts
├── public/                # Static assets
├── index.html             # HTML template
├── package.json           # Dependencies
├── tsconfig.json          # TypeScript config
├── vite.config.ts         # Vite config
├── tailwind.config.js     # Tailwind config
└── Dockerfile             # Container definition
```

### Key Components

#### Pages

1. **UploadPage**: File upload with drag-and-drop
2. **HostsPage**: List of all monitored hosts
3. **SnapshotsPage**: Snapshot history for a host
4. **DiffPage**: Visual comparison of two snapshots

#### API Client (`api/api.ts`)

Axios-based HTTP client with endpoints:
- `uploadSnapshot()`: Upload file
- `listHosts()`: Get all hosts
- `listSnapshotsByIP()`: Get host snapshots
- `compareSnapshots()`: Get diff report

#### Type Definitions (`types/index.ts`)

TypeScript interfaces matching backend models:
- `Snapshot`, `Service`, `Software`, `TLS`
- `DiffReport`, `PortChange`, `ServiceChange`

## Data Flow

### Upload Flow

```
User selects file
      │
      ▼
UploadPage validates file
      │
      ▼
FormData sent to /api/snapshots/upload
      │
      ▼
SnapshotHandler receives file
      │
      ├─▶ Validate file type (.json)
      ├─▶ Parse JSON
      ├─▶ Validate required fields
      │
      ▼
SnapshotService stores in DB
      │
      ├─▶ Serialize to JSON
      ├─▶ INSERT into snapshots table
      │
      ▼
Success response with IP & timestamp
      │
      ▼
UI shows success message
```

### Comparison Flow

```
User selects two snapshots
      │
      ▼
SnapshotsPage enables "Compare" button
      │
      ▼
Navigate to /diff/:oldId/:newId
      │
      ▼
DiffPage requests GET /api/diff/:id1/:id2
      │
      ▼
DiffHandler retrieves both snapshots
      │
      ├─▶ Validate same host
      │
      ▼
DiffService compares snapshots
      │
      ├─▶ Find added ports
      ├─▶ Find removed ports
      ├─▶ Compare existing services
      │   ├─▶ Status changes
      │   ├─▶ Software changes
      │   ├─▶ TLS changes
      │   └─▶ Vulnerability changes
      │
      ▼
Return DiffReport JSON
      │
      ▼
DiffPage renders visual comparison
      │
      ├─▶ Summary cards (added/removed/modified)
      ├─▶ Added ports section
      ├─▶ Removed ports section
      └─▶ Modified services section
```

## Deployment Architecture

### Docker Compose Setup

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Host                          │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │          scandiff-network (bridge)             │    │
│  │                                                 │    │
│  │  ┌──────────────────┐  ┌──────────────────┐  │    │
│  │  │   frontend       │  │     backend      │  │    │
│  │  │  (nginx:alpine)  │  │   (golang:1.21)  │  │    │
│  │  │   Port 3000      │──│    Port 8080     │  │    │
│  │  └──────────────────┘  └─────────┬────────┘  │    │
│  │                                   │            │    │
│  │                          ┌────────▼────────┐  │    │
│  │                          │   SQLite DB     │  │    │
│  │                          │  (volume mount) │  │    │
│  │                          └─────────────────┘  │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  Volumes:                                               │
│  - ./data:/app/data (SQLite persistence)               │
└─────────────────────────────────────────────────────────┘
```

### Container Details

**Frontend Container**
- Base: `nginx:alpine`
- Build: Multi-stage (Node.js build → Nginx serve)
- Exposes: Port 3000
- Proxies `/api/*` to backend

**Backend Container**
- Base: `golang:1.21-alpine`
- Build: Multi-stage (Go build → Alpine runtime)
- Exposes: Port 8080
- Mounts: `./data` for SQLite file

## Security Considerations

1. **Input Validation**: All JSON uploads are validated
2. **SQL Injection**: Parameterized queries only
3. **CORS**: Configured for specific origins
4. **Error Messages**: No sensitive data in error responses
5. **File Upload**: Limited to JSON files only

## Performance Considerations

1. **Database Indexing**: Indexes on `ip` and `timestamp`
2. **Frontend Bundling**: Vite optimizes build
3. **Nginx Caching**: Static assets cached
4. **Docker Multi-stage**: Minimal image sizes
5. **Connection Pooling**: Go database/sql handles pooling

## Scalability Considerations

### Current Limitations (SQLite)
- Single-writer constraint
- File-based storage
- Limited concurrent writes

### Future Scaling Options
1. **Database**: Migrate to PostgreSQL
2. **Caching**: Add Redis for frequently accessed data
3. **Load Balancing**: Multiple backend instances
4. **Object Storage**: Store snapshots in S3/MinIO/CloudFlare R2

## Testing Architecture

### Backend Tests
- **Unit Tests**: Service logic (diff algorithms)
- **Integration Tests**: API handlers with in-memory DB
- **Test Coverage**: Critical paths covered

### Frontend Tests
- **E2E Tests**: Playwright for user flows

## Monitoring & Observability

### Current Implementation
- Health check endpoint: `GET /health`
- Docker health checks
- Application logs