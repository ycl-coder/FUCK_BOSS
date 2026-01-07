// Application configuration

export const config = {
  // gRPC server URL
  // 使用相对路径，通过 Vite proxy 或 Nginx proxy 转发到后端
  // 开发环境：Vite proxy 会转发 /api 到 localhost:50051
  // 生产环境：Nginx 会转发 /api 到后端
  grpcUrl: import.meta.env.VITE_GRPC_URL || '',
  
  // API base URL
  // 使用相对路径，通过代理转发
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '',
  
  // Default pagination
  defaultPageSize: 20,
  
  // App info
  appName: 'Fuck Boss',
  appVersion: '1.0.0',
} as const

