// Application configuration from environment variables

const getEnvVar = (key: string, defaultValue: string = ''): string => {
  if (typeof window === 'undefined') {
    // Server-side
    return process.env[key] || defaultValue;
  }
  // Client-side - only NEXT_PUBLIC_ variables are available
  return (process.env[key] || defaultValue);
};

export const config = {
  api: {
    url: getEnvVar('NEXT_PUBLIC_API_URL', 'http://localhost:8080'),
    timeout: parseInt(getEnvVar('NEXT_PUBLIC_API_TIMEOUT', '30000')),
  },
  minio: {
    url: getEnvVar('NEXT_PUBLIC_MINIO_URL', 'http://localhost:9000'),
    consoleUrl: getEnvVar('NEXT_PUBLIC_MINIO_CONSOLE_URL', 'http://localhost:9001'),
  },
  app: {
    name: getEnvVar('NEXT_PUBLIC_APP_NAME', 'JuiceFS File Browser'),
    maxUploadSize: parseInt(getEnvVar('NEXT_PUBLIC_MAX_UPLOAD_SIZE', '104857600')),
    enableDelete: getEnvVar('NEXT_PUBLIC_ENABLE_DELETE', 'true') === 'true',
    enableUpload: getEnvVar('NEXT_PUBLIC_ENABLE_UPLOAD', 'true') === 'true',
  },
  ui: {
    defaultView: getEnvVar('NEXT_PUBLIC_DEFAULT_VIEW', 'tree') as 'tree' | 'flat',
    showHiddenFiles: getEnvVar('NEXT_PUBLIC_SHOW_HIDDEN_FILES', 'false') === 'true',
    filePreviewEnabled: getEnvVar('NEXT_PUBLIC_FILE_PREVIEW_ENABLED', 'true') === 'true',
  }
};

// Helper function to build API URLs
export const buildApiUrl = (path: string): string => {
  const baseUrl = config.api.url;
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  return `${baseUrl}/api${cleanPath}`;
};