#!/bin/bash
set -e

echo "JuiceFS Configuration:"
echo "  Meta: ${JUICEFS_META}"
echo "  Storage: ${JUICEFS_STORAGE}"
echo "  Bucket: ${JUICEFS_BUCKET}"

# Format JuiceFS if needed (only first time)
if [ "$JUICEFS_META" ] && [ "$JUICEFS_STORAGE" ]; then
    echo "Formatting JuiceFS (if not exists)..."
    juicefs format \
        --storage ${JUICEFS_STORAGE} \
        --bucket ${JUICEFS_BUCKET} \
        --access-key ${JUICEFS_ACCESS_KEY} \
        --secret-key ${JUICEFS_SECRET_KEY} \
        ${JUICEFS_META} \
        ${JUICEFS_NAME} || echo "Already formatted"
fi

# Mount JuiceFS
echo "Mounting JuiceFS to /jfs..."
juicefs mount ${JUICEFS_META} /jfs -d --no-usage-report --cache-size 100 &

# Wait for mount
echo "Waiting for mount..."
for i in {1..30}; do
    if mountpoint -q /jfs; then
        echo "JuiceFS mounted successfully"
        ls -la /jfs
        break
    fi
    sleep 1
done

if ! mountpoint -q /jfs; then
    echo "ERROR: JuiceFS mount failed after 30 seconds"
    exit 1
fi

# Start the Go application
echo "Starting display-sharing service on port ${PORT:-8090}..."
exec /app/main