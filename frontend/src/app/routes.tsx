import { Routes, Route } from 'react-router-dom'
import { Layout } from 'antd'
import { AppHeader } from '@/shared/components/AppHeader'
import { HomePage } from '@/features/post/pages/HomePage'
import { CreatePostPage } from '@/features/post/pages/CreatePostPage'
import { PostDetailPage } from '@/features/post/pages/PostDetailPage'
import { SearchPage } from '@/features/search/pages/SearchPage'
import './routes.css'

const { Content, Footer } = Layout

export function AppRoutes() {
  return (
    <Layout className="app-layout">
      <AppHeader />
      <Content className="app-content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/create" element={<CreatePostPage />} />
          <Route path="/post/:id" element={<PostDetailPage />} />
          <Route path="/search" element={<SearchPage />} />
        </Routes>
      </Content>
      <Footer className="app-footer">
        <div className="footer-content">
          <p>Fuck Boss © 2026</p>
          <p className="footer-subtitle">公司曝光平台 - 让职场更透明</p>
        </div>
      </Footer>
    </Layout>
  )
}

