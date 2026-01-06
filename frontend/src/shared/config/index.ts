// Application configuration

export const config = {
  // gRPC server URL
  grpcUrl: import.meta.env.VITE_GRPC_URL || 'http://localhost:50051',
  
  // API base URL
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:50051',
  
  // Default pagination
  defaultPageSize: 20,
  
  // App info
  appName: 'Fuck Boss',
  appVersion: '1.0.0',
} as const

