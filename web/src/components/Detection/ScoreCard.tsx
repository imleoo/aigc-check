import { Card, Row, Col, Tag, Statistic } from 'antd'
import ReactECharts from 'echarts-for-react'
import type { Score, MultimodalResult } from '../../types/detection'

interface ScoreCardProps {
  score: Score
  riskLevel: string
  multimodal?: MultimodalResult
}

const getRiskLevelColor = (level: string) => {
  switch (level) {
    case 'low': return 'green'
    case 'medium': return 'orange'
    case 'high': return 'red'
    case 'very_high': return 'red'
    default: return 'default'
  }
}

const getRiskLevelText = (level: string) => {
  switch (level) {
    case 'low': return '低风险'
    case 'medium': return '中等风险'
    case 'high': return '高风险'
    case 'very_high': return '极高风险'
    default: return level
  }
}

function ScoreCard({ score, riskLevel }: ScoreCardProps) {
  const radarOption = {
    title: { text: '五维度评分' },
    radar: {
      indicator: [
        { name: '词汇', max: 100 },
        { name: '句子', max: 100 },
        { name: '个性化', max: 100 },
        { name: '逻辑', max: 100 },
        { name: '情感', max: 100 }
      ]
    },
    series: [{
      type: 'radar',
      data: [{
        value: [
          score.dimensions?.vocabulary || 0,
          score.dimensions?.sentence || 0,
          score.dimensions?.personalization || 0,
          score.dimensions?.logic || 0,
          score.dimensions?.emotion || 0
        ]
      }]
    }]
  }

  return (
    <Card title="检测结果" style={{ marginBottom: 24 }}>
      <Row gutter={16}>
        <Col span={8}>
          <Statistic
            title="总分"
            value={score.total}
            suffix="/ 100"
            valueStyle={{ color: score.total > 75 ? '#3f8600' : '#cf1322' }}
          />
          <Tag color={getRiskLevelColor(riskLevel)} style={{ marginTop: 8 }}>
            {getRiskLevelText(riskLevel)}
          </Tag>
        </Col>
        <Col span={16}>
          <ReactECharts option={radarOption} style={{ height: 300 }} />
        </Col>
      </Row>
    </Card>
  )
}

export default ScoreCard
