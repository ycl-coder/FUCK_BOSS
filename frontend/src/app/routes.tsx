import { Routes, Route } from 'react-router-dom'
import { Layout } from 'antd'
import { AppHeader } from '@/shared/components/AppHeader'
import { HomePage } from '@/features/post/pages/HomePage'
import { CreatePostPage } from '@/features/post/pages/CreatePostPage'
import { PostDetailPage } from '@/features/post/pages/PostDetailPage'
import { SearchPage } from '@/features/search/pages/SearchPage'

const { Content, Footer } = Layout

export function AppRoutes() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <AppHeader />
      <Content style={{ padding: '24px', background: '#fff' }}>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/create" element={<CreatePostPage />} />
          <Route path="/post/:id" element={<PostDetailPage />} />
          <Route path="/search" element={<SearchPage />} />
        </Routes>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Fuck Boss ©2026 - 公司曝光平台
      </Footer>
    </Layout>
  )
}

