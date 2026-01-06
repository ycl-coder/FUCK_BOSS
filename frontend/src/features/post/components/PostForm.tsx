import { useState } from 'react'
import { Form, Input, Select, DatePicker, Button, message, Space } from 'antd'
import type { FormProps } from 'antd'
import dayjs, { type Dayjs } from 'dayjs'
import { useNavigate } from 'react-router-dom'
import { CITIES } from '@/shared/constants/cities'
import type { CreatePostRequest } from '@/shared/types'
import { contentServiceClient } from '@/api/grpc/contentClient'

const { TextArea } = Input

interface PostFormValues {
  company: string
  cityCode: string
  content: string
  occurredAt?: Dayjs
}

interface PostFormProps {
  onSuccess?: () => void
}

export function PostForm({ onSuccess }: PostFormProps) {
  const [form] = Form.useForm<PostFormValues>()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)

  const handleSubmit: FormProps<PostFormValues>['onFinish'] = async (values) => {
    setLoading(true)
    try {
      // Get city name from code
      const city = CITIES.find((c) => c.code === values.cityCode)
      if (!city) {
        message.error('请选择有效的城市')
        return
      }

      // Prepare request
      const request: CreatePostRequest = {
        company: values.company.trim(),
        cityCode: values.cityCode,
        cityName: city.name,
        content: values.content.trim(),
        occurredAt: values.occurredAt
          ? Math.floor(values.occurredAt.valueOf() / 1000) // Convert to Unix timestamp (seconds)
          : undefined,
      }

      // Call API
      const response = await contentServiceClient.createPost(request)

      message.success('发布成功！')
      form.resetFields()

      // Call success callback or navigate
      if (onSuccess) {
        onSuccess()
      } else {
        // Navigate to post detail page
        navigate(`/post/${response.postId}`)
      }
    } catch (error) {
      console.error('Failed to create post:', error)
      message.error(
        error instanceof Error ? error.message : '发布失败，请稍后重试'
      )
    } finally {
      setLoading(false)
    }
  }

  const handleReset = () => {
    form.resetFields()
  }

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={handleSubmit}
      autoComplete="off"
      size="large"
    >
      <Form.Item
        label="公司名称"
        name="company"
        rules={[
          { required: true, message: '请输入公司名称' },
          { min: 1, max: 100, message: '公司名称长度应在 1-100 个字符之间' },
        ]}
      >
        <Input placeholder="请输入公司名称" />
      </Form.Item>

      <Form.Item
        label="所在城市"
        name="cityCode"
        rules={[{ required: true, message: '请选择所在城市' }]}
      >
        <Select placeholder="请选择城市" showSearch optionFilterProp="label">
          {CITIES.map((city) => (
            <Select.Option key={city.code} value={city.code} label={city.name}>
              {city.name}
            </Select.Option>
          ))}
        </Select>
      </Form.Item>

      <Form.Item
        label="曝光内容"
        name="content"
        rules={[
          { required: true, message: '请输入曝光内容' },
          { min: 10, max: 5000, message: '内容长度应在 10-5000 个字符之间' },
        ]}
      >
        <TextArea
          placeholder="请详细描述公司的不当行为（至少 10 个字符）"
          rows={6}
          showCount
          maxLength={5000}
        />
      </Form.Item>

      <Form.Item
        label="发生时间（可选）"
        name="occurredAt"
        tooltip="如果不填写，将使用当前时间"
      >
        <DatePicker
          style={{ width: '100%' }}
          placeholder="选择发生时间"
          showTime
          format="YYYY-MM-DD HH:mm:ss"
          disabledDate={(current) => current && current > dayjs().endOf('day')}
        />
      </Form.Item>

      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={loading}>
            发布
          </Button>
          <Button onClick={handleReset} disabled={loading}>
            重置
          </Button>
        </Space>
      </Form.Item>
    </Form>
  )
}

