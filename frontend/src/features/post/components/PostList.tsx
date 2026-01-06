import { useState, useEffect } from 'react'
import { List, Card, Pagination, Empty, Spin, message, Tag, Typography, Space } from 'antd'
import type { Post } from '@/shared/types'
import { contentServiceClient } from '@/api/grpc/contentClient'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import './PostList.css'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const { Text, Paragraph } = Typography

interface PostListProps {
  cityCode?: string
  initialPage?: number
  pageSize?: number
  onPostClick?: (postId: string) => void
}

export function PostList({
  cityCode = '',
  initialPage = 1,
  pageSize = 20,
  onPostClick,
}: PostListProps) {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(false)
  const [currentPage, setCurrentPage] = useState(initialPage)
  const [total, setTotal] = useState(0)

  const loadPosts = async (page: number) => {
    setLoading(true)
    try {
      const response = await contentServiceClient.listPosts(cityCode, page, pageSize)
      setPosts(response.posts)
      setTotal(response.total)
      setCurrentPage(response.page)
    } catch (error) {
      console.error('Failed to load posts:', error)
      message.error(
        error instanceof Error ? error.message : '加载帖子列表失败，请稍后重试'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadPosts(currentPage)
  }, [cityCode, currentPage, pageSize])

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    // Scroll to top when page changes
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handlePostClick = (postId: string) => {
    if (onPostClick) {
      onPostClick(postId)
    }
  }

  if (loading && posts.length === 0) {
    return (
      <div className="post-list-loading">
        <Spin size="large" />
      </div>
    )
  }

  if (!loading && posts.length === 0) {
    return (
      <Empty
        description="暂无曝光内容"
        image={Empty.PRESENTED_IMAGE_SIMPLE}
      />
    )
  }

  return (
    <div>
      <List
        dataSource={posts}
        loading={loading}
        renderItem={(post) => (
          <List.Item className="post-list-item">
            <Card
              hoverable
              className="post-list-card"
              onClick={() => handlePostClick(post.id)}
            >
              <div className="post-list-header">
                <Space wrap>
                  <Text strong className="post-list-company">
                    {post.company}
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
                className="post-list-content"
              >
                {post.content}
              </Paragraph>
              <div className="post-list-footer">
                <Text type="secondary" className="post-list-time">
                  {dayjs.unix(post.createdAt).fromNow()}
                </Text>
              </div>
            </Card>
          </List.Item>
        )}
      />
      {total > pageSize && (
        <div className="post-list-pagination">
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

