// gRPC Web client configuration and utilities

import { config } from '@/shared/config'

// gRPC Web client options
export interface GrpcClientOptions {
  host?: string
  port?: number
  useTls?: boolean
}

// Get gRPC Web service URL
export function getGrpcWebUrl(options?: GrpcClientOptions): string {
  if (options?.host && options?.port) {
    const protocol = options.useTls ? 'https' : 'http'
    return `${protocol}://${options.host}:${options.port}`
  }
  return config.grpcUrl
}

// Convert gRPC error to application error
export function handleGrpcError(error: unknown): Error {
  if (error instanceof Error) {
    return error
  }
  return new Error(String(error))
}

