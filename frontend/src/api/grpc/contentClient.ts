// Content Service gRPC Web client
// This is a placeholder that will be replaced with generated code from proto files
// For now, we'll create a mock client interface

import type { Post, PostListResponse, CreatePostRequest, SearchRequest, SearchResponse } from '@/shared/types'

// Content Service client interface
export interface ContentServiceClient {
  createPost(request: CreatePostRequest): Promise<{ postId: string; createdAt: number }>
  listPosts(cityCode: string, page: number, pageSize: number): Promise<PostListResponse>
  getPost(postId: string): Promise<Post>
  searchPosts(request: SearchRequest): Promise<SearchResponse>
}

// Mock implementation (will be replaced with actual gRPC Web client)
class MockContentServiceClient implements ContentServiceClient {
  constructor(_baseUrl?: string) {
    // baseUrl will be used when implementing actual gRPC Web client
  }

  async createPost(_request: CreatePostRequest): Promise<{ postId: string; createdAt: number }> {
    // TODO: Replace with actual gRPC Web call
    // For now, return mock response
    throw new Error('gRPC Web client not implemented yet. Please generate TypeScript code from proto files.')
  }

  async listPosts(_cityCode: string, _page: number, _pageSize: number): Promise<PostListResponse> {
    // TODO: Replace with actual gRPC Web call
    throw new Error('gRPC Web client not implemented yet. Please generate TypeScript code from proto files.')
  }

  async getPost(_postId: string): Promise<Post> {
    // TODO: Replace with actual gRPC Web call
    throw new Error('gRPC Web client not implemented yet. Please generate TypeScript code from proto files.')
  }

  async searchPosts(_request: SearchRequest): Promise<SearchResponse> {
    // TODO: Replace with actual gRPC Web call
    throw new Error('gRPC Web client not implemented yet. Please generate TypeScript code from proto files.')
  }
}

// Create Content Service client instance
export function createContentServiceClient(baseUrl?: string): ContentServiceClient {
  return new MockContentServiceClient(baseUrl)
}

// Default client instance
export const contentServiceClient = createContentServiceClient()

