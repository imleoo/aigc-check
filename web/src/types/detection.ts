export interface DetectionRequest {
  text: string
  options: DetectionOptions
}

export interface DetectionOptions {
  enable_multimodal?: boolean
  enable_statistics?: boolean
  enable_semantic?: boolean
  language?: string
}

export interface DetectionResult {
  id: string
  request_id: string
  text: string
  score: Score
  risk_level: string
  rule_results: RuleResult[]
  suggestions: Suggestion[]
  multimodal?: MultimodalResult
  process_time: string
  detected_at: string
}

export interface Score {
  total: number
  dimensions?: {
    vocabulary: number
    sentence: number
    personalization: number
    logic: number
    emotion: number
  }
}

export interface RuleResult {
  rule_type: string
  rule_name: string
  description: string
  detected: boolean
  score: number
  severity: string
  matches: Match[]
  count: number
  threshold: number
  message: string
}

export interface Match {
  text: string
  position: number
  context: string
}

export interface Suggestion {
  id: string
  type: string
  severity: string
  message: string
  original_text?: string
  suggested_text?: string
  position?: number
}

export interface MultimodalResult {
  rule_layer_score: number
  statistics_layer_score: number
  semantic_layer_score: number
  final_score: number
  confidence: number
  detection_mode: string
}

export interface HistoryItem {
  id: string
  request_id: string
  text_preview: string
  score: number
  risk_level: string
  created_at: string
}

export interface HistoryListResult {
  total: number
  page: number
  page_size: number
  items: HistoryItem[]
}

export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}
