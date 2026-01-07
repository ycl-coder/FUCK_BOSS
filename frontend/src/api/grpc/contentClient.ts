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
    // Remove trailing slash
    this.baseUrl = (baseUrl || config.grpcUrl).replace(/\/$/, '')
  }

  private async call<TRequest, TResponse>(
    service: string,
    method: string,
    request: TRequest
  ): Promise<TResponse> {
    // gRPC Web URL format: /package.service/method
    const url = `${this.baseUrl}/${service}/${method}`
    
    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/grpc-web+json',
          'Accept': 'application/grpc-web+json',
          'X-Grpc-Web': '1',
        },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}: ${response.statusText}`
        try {
          const errorData = await response.json()
          errorMessage = errorData.message || errorData.error || errorMessage
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
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error(String(error))
    }
  }

  async createPost(request: CreatePostRequest): Promise<{ postId: string; createdAt: number }> {
    const response = await this.call<CreatePostRequest, { postId: string; createdAt: number }>(
      'content.v1.ContentService',
      'CreatePost',
      request
    )
    return response
  }

  async listPosts(cityCode: string, page: number, pageSize: number): Promise<PostListResponse> {
    const response = await this.call<
      { cityCode: string; page: number; pageSize: number },
      PostListResponse
    >(
      'content.v1.ContentService',
      'ListPosts',
      { cityCode: cityCode || '', page, pageSize }
    )
    return response
  }

  async getPost(postId: string): Promise<Post> {
    const response = await this.call<{ postId: string }, { post: Post }>(
      'content.v1.ContentService',
      'GetPost',
      { postId }
    )
    if (!response.post) {
      throw new Error('Post not found')
    }
    return response.post
  }

  async searchPosts(request: SearchRequest): Promise<SearchResponse> {
    const response = await this.call<SearchRequest, SearchResponse>(
      'content.v1.ContentService',
      'SearchPosts',
      {
        keyword: request.keyword,
        cityCode: request.cityCode || '',
        page: request.page || 1,
        pageSize: request.pageSize || 20,
      }
    )
    return response
  }
}

// Create Content Service client instance
export function createContentServiceClient(baseUrl?: string): ContentServiceClient {
  return new GrpcWebContentServiceClient(baseUrl)
}

// Default client instance
export const contentServiceClient = createContentServiceClient()
