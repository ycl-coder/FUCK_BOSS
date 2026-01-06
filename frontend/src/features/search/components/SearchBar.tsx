import { useState } from 'react'
import { Input, Select, Button, Space } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { CITIES } from '@/shared/constants/cities'

const { Search } = Input

interface SearchBarProps {
  onSearch: (keyword: string, cityCode?: string) => void
  loading?: boolean
  initialKeyword?: string
  initialCityCode?: string
}

export function SearchBar({
  onSearch,
  loading = false,
  initialKeyword = '',
  initialCityCode,
}: SearchBarProps) {
  const [keyword, setKeyword] = useState(initialKeyword)
  const [cityCode, setCityCode] = useState<string | undefined>(initialCityCode)

  const handleSearch = () => {
    if (keyword.trim()) {
      onSearch(keyword.trim(), cityCode)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  return (
    <Space.Compact style={{ width: '100%' }} size="large">
      <Search
        placeholder="输入关键词搜索（公司名称或内容）"
        value={keyword}
        onChange={(e) => setKeyword(e.target.value)}
        onKeyPress={handleKeyPress}
        onSearch={handleSearch}
        enterButton={<Button icon={<SearchOutlined />} loading={loading}>搜索</Button>}
        size="large"
        allowClear
      />
      <Select
        placeholder="选择城市（可选）"
        style={{ width: 150 }}
        value={cityCode}
        onChange={setCityCode}
        allowClear
      >
        {CITIES.map((city) => (
          <Select.Option key={city.code} value={city.code}>
            {city.name}
          </Select.Option>
        ))}
      </Select>
    </Space.Compact>
  )
}

