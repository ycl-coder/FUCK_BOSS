import { useParams } from 'react-router-dom'
import { PostDetail } from '../components/PostDetail'

export function PostDetailPage() {
  const { id } = useParams<{ id: string }>()

  if (!id) {
    return <div>无效的帖子 ID</div>
  }

  return <PostDetail postId={id} />
}

