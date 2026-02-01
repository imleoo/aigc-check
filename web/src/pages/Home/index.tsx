import { useState } from 'react'
import { Layout, message } from 'antd'
import TextInput from '../../components/Detection/TextInput'
import ResultPanel from '../../components/Detection/ResultPanel'
import { detect } from '../../api/detection'
import type { DetectionOptions, DetectionResult } from '../../types/detection'

const { Content } = Layout

function Home() {
  const [text, setText] = useState('')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<DetectionResult | null>(null)
  const [options, setOptions] = useState<DetectionOptions>({
    enable_multimodal: false,
    enable_statistics: true,
    enable_semantic: false,
    language: 'zh'
  })

  const handleDetect = async () => {
    if (!text.trim()) {
      message.warning('请输入待检测的文本')
      return
    }

    setLoading(true)
    try {
      const result = await detect({ text, options })
      setResult(result)
      message.success('检测完成')
    } catch (error) {
      console.error('检测失败:', error)
      message.error('检测失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Layout style={{ minHeight: '100vh', background: '#f0f2f5' }}>
      <Content style={{ padding: '24px', maxWidth: 1200, margin: '0 auto', width: '100%' }}>
        <h1 style={{ textAlign: 'center', marginBottom: 32 }}>
          AIGC 内容检测系统
        </h1>

        <TextInput
          value={text}
          onChange={setText}
          onDetect={handleDetect}
          loading={loading}
          options={options}
          onOptionsChange={setOptions}
        />

        <ResultPanel result={result} loading={loading} />
      </Content>
    </Layout>
  )
}

export default Home
