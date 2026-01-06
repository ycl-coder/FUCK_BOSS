import { useState } from 'react'
import { Card, Select, Space, Typography } from 'antd'
import { useNavigate } from 'react-router-dom'
import { PostList } from '../components/PostList'
import { CITIES } from '@/shared/constants/cities'
import './HomePage.css'

const { Title } = Typography

export function HomePage() {
  const navigate = useNavigate()
  const [cityCode, setCityCode] = useState<string>('')

  const handlePostClick = (postId: string) => {
    navigate(`/post/${postId}`)
  }

  return (
    <div className="home-page">
      <Card className="home-page-card">
        <Space direction="vertical" className="home-page-content" size="large">
          <div className="home-page-header">
            <Title level={2} className="home-page-title">
              公司曝光平台
            </Title>
            <p className="home-page-subtitle">
              匿名分享职场经历，让信息更透明
            </p>
            <Select
              placeholder="选择城市筛选（留空显示全部）"
              className="home-page-filter"
              allowClear
              value={cityCode || undefined}
              onChange={(value) => setCityCode(value || '')}
              size="large"
            >
              {CITIES.map((city) => (
                <Select.Option key={city.code} value={city.code}>
                  {city.name}
                </Select.Option>
              ))}
            </Select>
          </div>
          <PostList cityCode={cityCode} onPostClick={handlePostClick} />
        </Space>
      </Card>
    </div>
  )
}

