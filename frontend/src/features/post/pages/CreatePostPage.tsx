import { Card, Typography } from 'antd'
import { PostForm } from '../components/PostForm'
import './CreatePostPage.css'

const { Title } = Typography

export function CreatePostPage() {
  return (
    <div className="create-post-page">
      <Card className="create-post-card">
        <div className="create-post-header">
          <Title level={2} className="create-post-title">
            发布曝光
          </Title>
          <p className="create-post-subtitle">
            分享您的职场经历，帮助更多人了解真实的工作环境
          </p>
        </div>
        <PostForm />
      </Card>
    </div>
  )
}

