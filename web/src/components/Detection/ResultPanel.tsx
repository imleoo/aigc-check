import { Space } from 'antd'
import ScoreCard from './ScoreCard'
import RuleResultCard from './RuleResultCard'
import type { DetectionResult } from '../../types/detection'

interface ResultPanelProps {
  result: DetectionResult | null
  loading: boolean
}

function ResultPanel({ result, loading }: ResultPanelProps) {
  if (loading) {
    return <div>检测中...</div>
  }

  if (!result) {
    return null
  }

  return (
    <Space direction="vertical" style={{ width: '100%' }} size="large">
      <ScoreCard
        score={result.score}
        riskLevel={result.risk_level}
        multimodal={result.multimodal}
      />
      <RuleResultCard ruleResults={result.rule_results} />
    </Space>
  )
}

export default ResultPanel
