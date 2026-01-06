import { Link, useLocation } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { HomeOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons'

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
    <Header style={{ display: 'flex', alignItems: 'center' }}>
      <div style={{ color: '#fff', fontSize: '20px', fontWeight: 'bold', marginRight: '24px' }}>
        Fuck Boss
      </div>
      <Menu
        theme="dark"
        mode="horizontal"
        selectedKeys={[location.pathname]}
        items={menuItems}
        style={{ flex: 1, minWidth: 0 }}
      />
    </Header>
  )
}

