import { useState } from 'react'
import { Card, Select, Space, Typography, Input, Row, Col } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { PostList } from '../components/PostList'
import { CITIES } from '@/shared/constants/cities'
import './HomePage.css'

const { Title } = Typography
const { Search } = Input

export function HomePage() {
  const navigate = useNavigate()
  const [cityCode, setCityCode] = useState<string>('')
  const [searchKeyword, setSearchKeyword] = useState<string>('')

  const handlePostClick = (postId: string) => {
    navigate(`/post/${postId}`)
  }

  const handleSearch = (value: string) => {
    if (value.trim()) {
      navigate(`/search?keyword=${encodeURIComponent(value.trim())}${cityCode ? `&city=${cityCode}` : ''}`)
    }
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
            <Row gutter={16} className="home-page-filters">
              <Col xs={24} sm={14} md={16}>
                <Search
                  placeholder="搜索公司名称或曝光内容"
                  allowClear
                  enterButton={<SearchOutlined />}
                  size="large"
                  onSearch={handleSearch}
                  value={searchKeyword}
                  onChange={(e) => setSearchKeyword(e.target.value)}
                  className="home-page-search"
                />
              </Col>
              <Col xs={24} sm={10} md={8}>
                <Select
                  placeholder="选择城市筛选（留空显示全部）"
                  className="home-page-filter"
                  allowClear
                  value={cityCode || undefined}
                  onChange={(value) => setCityCode(value || '')}
                  size="large"
                  style={{ width: '100%' }}
                >
                  {CITIES.map((city) => (
                    <Select.Option key={city.code} value={city.code}>
                      {city.name}
                    </Select.Option>
                  ))}
                </Select>
              </Col>
            </Row>
          </div>
          <PostList cityCode={cityCode} onPostClick={handlePostClick} />
        </Space>
      </Card>
    </div>
  )
}

