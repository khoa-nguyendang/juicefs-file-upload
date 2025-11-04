# JuiceFS File Browser

A modern web-based file browser application that provides a user-friendly interface for managing files in MinIO S3-compatible object storage. Built with Go backend and React/Next.js frontend, this solution is designed to work seamlessly in both Docker and Kubernetes environments.

## Features

- **Tree View Navigation**: Browse files and folders in a hierarchical tree structure
- **Flat View Option**: Switch between tree and flat view modes
- **File Upload**: Upload individual files or entire folders while preserving directory structure
- **Drag & Drop**: Intuitive drag-and-drop interface for file uploads
- **Progress Tracking**: Real-time upload progress with visual feedback
- **File Operations**: Download and delete files directly from the UI
- **Folder Structure Preservation**: Maintains original folder hierarchy during uploads
- **S3-Compatible**: Uses MinIO for reliable, scalable object storage
- **Kubernetes Ready**: Designed for containerized deployments

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   React UI  │────▶│  Go Server  │────▶│    MinIO    │
│  (Next.js)  │     │    (API)    │     │  (Storage)  │
└─────────────┘     └─────────────┘     └─────────────┘
     :3000              :8080               :9000
```

## Quick Start

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/khoa-nguyendang/juicefs-file-upload.git
cd juicefs-file-upload

# Start all services
docker-compose up --build
# Or use the Makefile
make up

# Access the application
# UI: http://localhost:3000
# API: http://localhost:8080
# MinIO Console: http://localhost:9001

# Stop services
docker-compose down
# Or use the Makefile
make down
```

### Default Credentials

- MinIO Root User: `minioadmin`
- MinIO Root Password: `minioadmin123`
- MinIO Access Key: `juicefs`
- MinIO Secret Key: `juicefs123`

## Why MinIO Instead of JuiceFS FUSE?

While JuiceFS is excellent for distributed file systems, we chose direct MinIO integration for several reasons:

1. **Kubernetes Compatibility**: MinIO uses REST APIs instead of FUSE mounts, making it fully compatible with Kubernetes pod orchestration
2. **Simplicity**: Direct S3 API calls eliminate the complexity of filesystem mounts
3. **Scalability**: Both the Go server and UI can scale horizontally without shared filesystem dependencies
4. **Performance**: Direct object storage operations without filesystem abstraction layer

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/list?path=/` | List files in directory |
| POST | `/api/upload` | Upload file with optional path |
| GET | `/api/download/{filename}` | Download file |
| DELETE | `/api/delete/{filename}` | Delete file |

## Testing

### Test File Upload

```bash
# Upload a single file
curl -X POST -F "file=@test.txt" -F "path=/documents" \
  http://localhost:8080/api/upload

# List files
curl http://localhost:8080/api/list?path=/

# Download file
curl -O http://localhost:8080/api/download/test.txt
```

### Test Folder Upload

Use the provided test HTML file:

```bash
open test-folder-upload.html
# Select a folder and upload with structure preservation
```

## Production Deployment

### Kubernetes Deployment

See [KUBERNETES_DEPLOYMENT.md](KUBERNETES_DEPLOYMENT.md) for detailed Kubernetes deployment instructions.

### Building Images

```bash
# Build server image
docker build -t your-registry/file-browser-server:latest ./server

# Build UI image
docker build -t your-registry/file-browser-ui:latest ./ui

# Push to registry
docker push your-registry/file-browser-server:latest
docker push your-registry/file-browser-ui:latest
```

## Environment Variables

### Go Server
- `SERVER_PORT`: Server port (default: 8080)
- `MINIO_ENDPOINT`: MinIO endpoint
- `MINIO_ACCESS_KEY`: MinIO access key
- `MINIO_SECRET_KEY`: MinIO secret key
- `MINIO_USE_SSL`: Use SSL for MinIO (default: false)

### UI Application
- `NEXT_PUBLIC_API_URL`: Backend API URL
- `NEXT_PUBLIC_API_TIMEOUT`: API request timeout
- `NEXT_PUBLIC_MAX_UPLOAD_SIZE`: Maximum upload size in bytes
- `NEXT_PUBLIC_DEFAULT_VIEW`: Default view mode (tree/flat)
- `NEXT_PUBLIC_ENABLE_DELETE`: Enable delete functionality
- `NEXT_PUBLIC_ENABLE_UPLOAD`: Enable upload functionality

## License

MIT License - see LICENSE file for details