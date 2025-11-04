# Display Sharing Service - JuiceFS POSIX Demo

This service demonstrates the real value of JuiceFS by using POSIX filesystem operations to access files stored in MinIO S3.

## Architecture Comparison

### Original Service (Direct S3)
```
Browser → Go API → MinIO S3 API → Storage
         (HTTP)    (Object API)
```

### Display Sharing Service (JuiceFS POSIX)
```
Browser → Go API → JuiceFS Mount → MinIO S3 → Storage
         (HTTP)    (POSIX ops)     (Hidden)
```

## Key Differences

| Feature | Direct S3 Service | JuiceFS POSIX Service |
|---------|------------------|----------------------|
| File Access | S3 API calls | Standard Go file operations |
| Code | `minioClient.GetObject()` | `os.ReadFile()` |
| Directory Listing | S3 ListObjects | `os.ReadDir()` |
| File Delete | S3 DeleteObject | `os.Remove()` |
| Compatibility | S3-aware apps only | ANY application |
| File Locking | Not supported | POSIX file locks |
| Atomic Operations | Limited | Full POSIX semantics |

## API Endpoints

### List Files (POSIX `readdir`)
```bash
GET /api/list?path=/
# Uses: os.ReadDir("/jfs/path")
```

### View File Content (POSIX `read`)
```bash
GET /api/view?path=/file.txt
# Uses: os.ReadFile("/jfs/file.txt")
```

### Delete File (POSIX `unlink`)
```bash
DELETE /api/delete?path=/file.txt
# Uses: os.Remove("/jfs/file.txt")
```

### Download File (POSIX `open` + `read`)
```bash
GET /api/download?path=/file.pdf
# Uses: os.Open("/jfs/file.pdf") + io.Copy()
```

### Health Check
```bash
GET /api/health
# Verifies JuiceFS mount at /jfs
```

### Statistics
```bash
GET /api/stats
# Walks filesystem tree using filepath.WalkDir()
```