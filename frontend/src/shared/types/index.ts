// Shared types for the application

// Post related types
export interface Post {
  id: string
  company: string
  cityCode: string
  cityName: string
  content: string
  occurredAt?: number // Unix timestamp (optional)
  createdAt: number // Unix timestamp
}

export interface PostListResponse {
  posts: Post[]
  total: number
  page: number
  pageSize: number
}

export interface CreatePostRequest {
  company: string
  cityCode: string
  cityName: string
  content: string
  occurredAt?: number // Unix timestamp (optional)
}

// City related types
export interface City {
  code: string
  name: string
}

// Search related types
export interface SearchRequest {
  keyword: string
  cityCode?: string
  page?: number
  pageSize?: number
}

export interface SearchResponse {
  posts: Post[]
  total: number
  page: number
  pageSize: number
}

// API Error types
export interface ApiError {
  code: string
  message: string
  details?: Record<string, unknown>
}

