import { Link, useLocation } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { HomeOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons'
import './AppHeader.css'

const { Header } = Layout

export function AppHeader() {
  const location = useLocation()

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: <Link to="/">首页</Link>,
    },
    {
      key: '/create',
      icon: <PlusOutlined />,
      label: <Link to="/create">发布曝光</Link>,
    },
    {
      key: '/search',
      icon: <SearchOutlined />,
      label: <Link to="/search">搜索</Link>,
    },
  ]

  return (
    <Header className="app-header">
      <div className="app-header-brand">
        <Link to="/" className="brand-link">
          <span className="brand-name">Fuck Boss</span>
          <span className="brand-subtitle">公司曝光平台</span>
        </Link>
      </div>
      <Menu
        theme="dark"
        mode="horizontal"
        selectedKeys={[location.pathname]}
        items={menuItems}
        className="app-header-menu"
      />
    </Header>
  )
}

