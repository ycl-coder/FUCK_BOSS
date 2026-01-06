import { useState, useEffect } from 'react'
import { Card, Spin, Empty, message, Tag, Typography, Space, Button } from 'antd'
import { useNavigate } from 'react-router-dom'
import type { Post } from '@/shared/types'
import { contentServiceClient } from '@/api/grpc/contentClient'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import './PostDetail.css'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const { Title, Paragraph, Text } = Typography

interface PostDetailProps {
  postId: string
}

export function PostDetail({ postId }: PostDetailProps) {
  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    const loadPost = async () => {
      setLoading(true)
      try {
        const data = await contentServiceClient.getPost(postId)
        setPost(data)
      } catch (error) {
        console.error('Failed to load post:', error)
        const errorMessage =
          error instanceof Error ? error.message : '加载帖子详情失败'
        
        // Check if it's a 404 error
        if (errorMessage.includes('not found') || errorMessage.includes('NotFound')) {
          message.error('帖子不存在')
        } else {
          message.error(errorMessage)
        }
        setPost(null)
      } finally {
        setLoading(false)
      }
    }

    if (postId) {
      loadPost()
    }
  }, [postId])

  if (loading) {
    return (
      <div className="post-detail-loading">
        <Spin size="large" />
      </div>
    )
  }

  if (!post) {
    return (
      <div className="post-detail-empty">
        <Empty
          description="帖子不存在或已被删除"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        >
          <Button type="primary" onClick={() => navigate('/')}>
            返回首页
          </Button>
        </Empty>
      </div>
    )
  }

  return (
    <div className="post-detail-container">
      <Card className="post-detail-card">
        <Space direction="vertical" className="post-detail-content" size="large">
          <div>
            <Title level={2} className="post-detail-title">
              {post.company}
            </Title>
            <Space wrap>
              <Tag color="blue">{post.cityName}</Tag>
              {post.occurredAt && (
                <Tag color="orange">
                  发生时间：{dayjs.unix(post.occurredAt).format('YYYY-MM-DD HH:mm:ss')}
                </Tag>
              )}
              <Text type="secondary" className="post-detail-time">
                发布时间：{dayjs.unix(post.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Text>
            </Space>
          </div>

          <div>
            <Title level={4} className="post-detail-section-title">
              曝光内容
            </Title>
            <Paragraph className="post-detail-content-text">
              {post.content}
            </Paragraph>
          </div>

          <div className="post-detail-actions">
            <Space wrap>
              <Button onClick={() => navigate('/')}>返回首页</Button>
              <Button type="primary" onClick={() => navigate('/create')}>
                发布新曝光
              </Button>
            </Space>
          </div>
        </Space>
      </Card>
    </div>
  )
}

