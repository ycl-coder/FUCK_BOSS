import { useState, useEffect } from 'react'
import { Card, Typography } from 'antd'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { SearchBar } from '../components/SearchBar'
import { SearchResults } from '../components/SearchResults'
import './SearchPage.css'

const { Title } = Typography

export function SearchPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [keyword, setKeyword] = useState('')
  const [cityCode, setCityCode] = useState<string | undefined>()

  // 从 URL 参数初始化搜索关键词和城市
  useEffect(() => {
    const urlKeyword = searchParams.get('keyword')
    const urlCity = searchParams.get('city')
    if (urlKeyword) {
      setKeyword(urlKeyword)
    }
    if (urlCity) {
      setCityCode(urlCity)
    }
  }, [searchParams])

  const handleSearch = (searchKeyword: string, searchCityCode?: string) => {
    setKeyword(searchKeyword)
    setCityCode(searchCityCode)
    // 更新 URL 参数
    const params = new URLSearchParams()
    if (searchKeyword) {
      params.set('keyword', searchKeyword)
    }
    if (searchCityCode) {
      params.set('city', searchCityCode)
    }
    navigate(`/search?${params.toString()}`, { replace: true })
  }

  const handlePostClick = (postId: string) => {
    navigate(`/post/${postId}`)
  }

  return (
    <div className="search-page">
      <Card className="search-page-card">
        <Title level={2} className="search-page-title">
          搜索曝光内容
        </Title>
        <p className="search-page-subtitle">
          输入关键词搜索公司名称或曝光内容
        </p>
        <div className="search-page-bar">
          <SearchBar
            onSearch={handleSearch}
            initialKeyword={keyword}
            initialCityCode={cityCode}
          />
        </div>
        {keyword && (
          <SearchResults
            keyword={keyword}
            cityCode={cityCode}
            onPostClick={handlePostClick}
          />
        )}
      </Card>
    </div>
  )
}

