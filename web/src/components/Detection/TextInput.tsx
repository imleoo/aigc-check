import { Card, Input, Button, Space, Switch, Form } from 'antd'
import type { DetectionOptions } from '../../types/detection'

const { TextArea } = Input

interface TextInputProps {
  value: string
  onChange: (value: string) => void
  onDetect: () => void
  loading: boolean
  options: DetectionOptions
  onOptionsChange: (options: DetectionOptions) => void
}

function TextInput({
  value,
  onChange,
  onDetect,
  loading,
  options,
  onOptionsChange
}: TextInputProps) {
  const wordCount = value.length

  return (
    <Card title="文本输入" style={{ marginBottom: 24 }}>
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        <TextArea
          value={value}
          onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => onChange(e.target.value)}
          placeholder="请输入待检测的文本内容..."
          rows={10}
          showCount
          maxLength={10000}
        />

        <Form layout="inline">
          <Form.Item label="多模态检测">
            <Switch
              checked={options.enable_multimodal}
              onChange={(checked: boolean) =>
                onOptionsChange({ ...options, enable_multimodal: checked })
              }
            />
          </Form.Item>
          <Form.Item label="统计分析">
            <Switch
              checked={options.enable_statistics}
              onChange={(checked: boolean) =>
                onOptionsChange({ ...options, enable_statistics: checked })
              }
            />
          </Form.Item>
        </Form>

        <Space>
          <Button
            type="primary"
            size="large"
            onClick={onDetect}
            loading={loading}
            disabled={!value.trim()}
          >
            开始检测
          </Button>
          <span style={{ color: '#999' }}>
            字数: {wordCount} / 10000
          </span>
        </Space>
      </Space>
    </Card>
  )
}

export default TextInput
