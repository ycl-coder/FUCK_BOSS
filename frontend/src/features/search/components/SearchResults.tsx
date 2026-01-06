import { useState, useEffect, useMemo } from 'react'
import { List, Card, Pagination, Empty, Spin, Tag, Typography, Space, message } from 'antd'
import type { Post } from '@/shared/types'
import { contentServiceClient } from '@/api/grpc/contentClient'
import { CITIES } from '@/shared/constants/cities'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const { Text, Paragraph } = Typography

interface SearchResultsProps {
  keyword: string
  cityCode?: string
  onPostClick?: (postId: string) => void
  pageSize?: number
}

// Highlight keyword in text
function highlightText(text: string, keyword: string): React.ReactNode {
  if (!keyword) return text

  const parts = text.split(new RegExp(`(${keyword})`, 'gi'))
  return parts.map((part, index) =>
    part.toLowerCase() === keyword.toLowerCase() ? (
      <mark key={index} className="search-highlight">
        {part}
      </mark>
    ) : (
      part
    )
  )
}

export function SearchResults({
  keyword,
  cityCode,
  onPostClick,
  pageSize = 20,
}: SearchResultsProps) {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(false)
  const [currentPage, setCurrentPage] = useState(1)
  const [total, setTotal] = useState(0)

  const loadResults = async (page: number) => {
    if (!keyword.trim()) {
      setPosts([])
      setTotal(0)
      return
    }

    setLoading(true)
    try {
      const response = await contentServiceClient.searchPosts({
        keyword: keyword.trim(),
        cityCode,
        page,
        pageSize,
      })
      setPosts(response.posts)
      setTotal(response.total)
      setCurrentPage(response.page)
    } catch (error) {
      console.error('Failed to search posts:', error)
      message.error(
        error instanceof Error ? error.message : '搜索失败，请稍后重试'
      )
      setPosts([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    // Reset to page 1 when keyword or cityCode changes
    setCurrentPage(1)
    loadResults(1)
  }, [keyword, cityCode])

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    loadResults(page)
    // Scroll to top when page changes
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handlePostClick = (postId: string) => {
    if (onPostClick) {
      onPostClick(postId)
    }
  }

  // Memoize highlighted content
  const highlightedPosts = useMemo(() => {
    return posts.map((post) => ({
      ...post,
      highlightedCompany: highlightText(post.company, keyword),
      highlightedContent: highlightText(post.content, keyword),
    }))
  }, [posts, keyword])

  if (loading && posts.length === 0) {
    return (
      <div className="search-results-loading">
        <Spin size="large" />
      </div>
    )
  }

  if (!keyword.trim()) {
    return (
      <Empty
        description="请输入搜索关键词"
        image={Empty.PRESENTED_IMAGE_SIMPLE}
      />
    )
  }

  if (!loading && posts.length === 0) {
    return (
      <Empty
        description={`未找到包含"${keyword}"的曝光内容`}
        image={Empty.PRESENTED_IMAGE_SIMPLE}
      />
    )
  }

  return (
    <div>
      <div className="search-results-header">
        <Text type="secondary" className="search-results-count">
          找到 {total} 条相关结果
          {cityCode && (
            <span>
              {' '}
              （城市：{CITIES.find((c) => c.code === cityCode)?.name || cityCode}）
            </span>
          )}
        </Text>
      </div>

      <List
        dataSource={highlightedPosts}
        loading={loading}
        renderItem={(post) => (
          <List.Item className="search-results-item">
            <Card
              hoverable
              className="search-results-card"
              onClick={() => handlePostClick(post.id)}
            >
              <div className="search-results-header-info">
                <Space wrap>
                  <Text strong className="search-results-company">
                    {post.highlightedCompany}
                  </Text>
                  <Tag color="blue">{post.cityName}</Tag>
                  {post.occurredAt && (
                    <Tag color="orange">
                      {dayjs.unix(post.occurredAt).format('YYYY-MM-DD')}
                    </Tag>
                  )}
                </Space>
              </div>
              <Paragraph
                ellipsis={{ rows: 3, expandable: false }}
                className="search-results-content"
              >
                {post.highlightedContent}
              </Paragraph>
              <div className="search-results-footer">
                <Text type="secondary" className="search-results-time">
                  {dayjs.unix(post.createdAt).fromNow()}
                </Text>
              </div>
            </Card>
          </List.Item>
        )}
      />
      {total > pageSize && (
        <div className="search-results-pagination">
          <Pagination
            current={currentPage}
            total={total}
            pageSize={pageSize}
            onChange={handlePageChange}
            showSizeChanger={false}
            showQuickJumper
            showTotal={(total, range) =>
              `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
            }
          />
        </div>
      )}
    </div>
  )
}

