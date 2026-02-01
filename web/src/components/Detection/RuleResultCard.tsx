import { Card, List, Tag, Badge, Progress } from 'antd'
import type { RuleResult } from '../../types/detection'

interface RuleResultCardProps {
  ruleResults: RuleResult[]
}

function RuleResultCard({ ruleResults }: RuleResultCardProps) {
  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'red'
      case 'high': return 'orange'
      case 'medium': return 'gold'
      case 'low': return 'blue'
      default: return 'default'
    }
  }

  return (
    <Card title="规则匹配结果" style={{ marginBottom: 24 }}>
      <List
        dataSource={ruleResults}
        renderItem={(item: RuleResult) => (
          <List.Item>
            <List.Item.Meta
              title={
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <span>{item.rule_name}</span>
                  {item.detected && (
                    <Tag color="red">检测到</Tag>
                  )}
                  <Tag color={getSeverityColor(item.severity)}>
                    {item.severity}
                  </Tag>
                </div>
              }
              description={
                <div>
                  <div>{item.description}</div>
                  <div style={{ marginTop: 8 }}>
                    <Progress
                      percent={item.score}
                      size="small"
                      status={item.score >= 80 ? 'success' : item.score >= 60 ? 'normal' : 'exception'}
                    />
                  </div>
                  <div style={{ marginTop: 4, fontSize: 12, color: '#666' }}>
                    {item.message}
                  </div>
                </div>
              }
            />
            {item.matches && item.matches.length > 0 && (
              <Badge count={item.matches.length} />
            )}
          </List.Item>
        )}
      />
    </Card>
  )
}

export default RuleResultCard
