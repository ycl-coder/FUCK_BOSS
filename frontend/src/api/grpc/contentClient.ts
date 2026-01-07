// Content Service gRPC Web client implementation
// Uses fetch API with gRPC Web protocol to communicate with backend

import { config } from '@/shared/config'
import type { Post, PostListResponse, CreatePostRequest, SearchRequest, SearchResponse } from '@/shared/types'

// Content Service client interface
export interface ContentServiceClient {
  createPost(request: CreatePostRequest): Promise<{ postId: string; createdAt: number }>
  listPosts(cityCode: string, page: number, pageSize: number): Promise<PostListResponse>
  getPost(postId: string): Promise<Post>
  searchPosts(request: SearchRequest): Promise<SearchResponse>
}

// Real gRPC Web client implementation
// gRPC Web uses HTTP/1.1 with specific headers and URL format
class GrpcWebContentServiceClient implements ContentServiceClient {
  private baseUrl: string

  constructor(baseUrl?: string) {
    // Use provided baseUrl, or config, or empty string (relative path)
    // Empty string means use relative paths, which will go through Vite proxy or Nginx proxy
    const url = baseUrl || config.grpcUrl || ''
    this.baseUrl = url.replace(/\/$/, '')
  }

  private async call<TRequest, TResponse>(
    endpoint: string,
    request: TRequest,
    method: 'GET' | 'POST' = 'POST'
  ): Promise<TResponse> {
    // Use REST API endpoint
    // If baseUrl is empty, use relative path (will use current domain)
    // This allows Vite proxy (dev) or Nginx proxy (prod) to forward to backend
    const url = this.baseUrl ? `${this.baseUrl}${endpoint}` : endpoint
    
    try {
      const options: RequestInit = {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      }

      if (method === 'POST' && request) {
        options.body = JSON.stringify(request)
      } else if (method === 'GET' && request) {
        // Convert request to query parameters
        const params = new URLSearchParams()
        Object.entries(request as Record<string, unknown>).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== '') {
            params.append(key, String(value))
          }
        })
        const queryString = params.toString()
        if (queryString) {
          const separator = url.includes('?') ? '&' : '?'
          const urlWithQuery = `${url}${separator}${queryString}`
          const response = await fetch(urlWithQuery, options)
          return this.handleResponse<TResponse>(response)
        }
      }

      const response = await fetch(url, options)
      return this.handleResponse<TResponse>(response)
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error(String(error))
    }
  }

  private async handleResponse<TResponse>(response: Response): Promise<TResponse> {
    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`
      try {
        const errorData = await response.json()
        errorMessage = errorData.error || errorData.message || errorMessage
      } catch {
        const errorText = await response.text()
        if (errorText) {
          errorMessage = errorText
        }
      }
      throw new Error(errorMessage)
    }

    const data = await response.json()
    return data as TResponse
  }

  async createPost(request: CreatePostRequest): Promise<{ postId: string; createdAt: number }> {
    return this.call<CreatePostRequest, { postId: string; createdAt: number }>(
      '/api/posts',
      request,
      'POST'
    )
  }

  async listPosts(cityCode: string, page: number, pageSize: number): Promise<PostListResponse> {
    return this.call<
      { cityCode: string; page: number; pageSize: number },
      PostListResponse
    >(
      '/api/posts',
      { cityCode: cityCode || '', page, pageSize },
      'GET'
    )
  }

  async getPost(postId: string): Promise<Post> {
    return this.call<never, Post>(
      `/api/posts/${postId}`,
      undefined as never,
      'GET'
    )
  }

  async searchPosts(request: SearchRequest): Promise<SearchResponse> {
    return this.call<SearchRequest, SearchResponse>(
      '/api/posts/search',
      {
        keyword: request.keyword,
        cityCode: request.cityCode,
        page: request.page || 1,
        pageSize: request.pageSize || 20,
      },
      'POST'
    )
  }
}

// Create Content Service client instance
export function createContentServiceClient(baseUrl?: string): ContentServiceClient {
  return new GrpcWebContentServiceClient(baseUrl)
}

// Default client instance
export const contentServiceClient = createContentServiceClient()
