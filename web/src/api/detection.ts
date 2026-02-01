import axios from 'axios'
import type {
  DetectionRequest,
  DetectionResult,
  ApiResponse,
  HistoryListResult
} from '../types/detection'

// 创建 axios 实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// 执行检测
export const detect = async (
  request: DetectionRequest
): Promise<DetectionResult> => {
  const response = await api.post<any, ApiResponse<DetectionResult>>(
    '/detect',
    request
  )
  return response.data!
}

// 根据 ID 获取检测结果
export const getDetectionById = async (
  id: string
): Promise<DetectionResult> => {
  const response = await api.get<any, ApiResponse<DetectionResult>>(
    `/detect/${id}`
  )
  return response.data!
}

// 获取历史记录列表
export const getHistory = async (
  page = 1,
  pageSize = 20,
  sort = 'created_at',
  order = 'desc'
): Promise<HistoryListResult> => {
  const response = await api.get<any, ApiResponse<HistoryListResult>>(
    '/history',
    {
      params: { page, page_size: pageSize, sort, order }
    }
  )
  return response.data!
}

// 根据 ID 获取历史记录
export const getHistoryById = async (
  id: string
): Promise<DetectionResult> => {
  const response = await api.get<any, ApiResponse<DetectionResult>>(
    `/history/${id}`
  )
  return response.data!
}

// 删除历史记录
export const deleteHistory = async (id: string): Promise<void> => {
  await api.delete(`/history/${id}`)
}

// 删除所有历史记录
export const deleteAllHistory = async (): Promise<void> => {
  await api.delete('/history')
}

export default api
