# Gemini API 集成指南

## 概述

本文档介绍如何在 AIGC-Check 中集成和使用 Google Gemini API 进行语义分析。

## 前置要求

### 1. 获取 Gemini API Key

访问 [Google AI Studio](https://makersuite.google.com/app/apikey) 获取 API Key。

### 2. 环境准备

确保系统已安装：
- Go 1.25+
- 网络连接（访问 Google API）

## 配置方式

### 方式 1: 环境变量（推荐）

```bash
# 设置环境变量
export GEMINI_API_KEY=your_api_key_here

# 验证设置
echo $GEMINI_API_KEY
```

### 方式 2: 命令行参数

```bash
aigc-check -f sample.txt -m -s -g --api-key YOUR_API_KEY
```

### 方式 3: 配置文件

编辑 `configs/aigc-check.yaml`:

```yaml
gemini:
  enabled: true
  api_key: "your_api_key_here"  # 不推荐，建议使用环境变量
  model: "gemini-pro"
  temperature: 0.3
  max_tokens: 500
  timeout: 30s
```

**注意**: 出于安全考虑，不建议在配置文件中直接写入 API Key。

## 使用方法

### 基础使用

```bash
# 设置 API Key
export GEMINI_API_KEY=your_api_key_here

# 启用完整多模态检测
aigc-check -f sample.txt -m -s -g

# 显示详细分析结果
aigc-check -f sample.txt -m -s -g --verbose
```

### 检测流程

当启用 Gemini API 时，系统会按照以下流程执行：

1. **规则检测层**: 快速识别明显 AI 特征
2. **置信度评估**: 计算规则检测置信度
3. **统计分析层**: 如果置信度 < 0.85，触发统计分析
4. **置信度评估**: 计算规则和统计的平均置信度
5. **语义分析层**: 如果平均置信度 < 0.6，调用 Gemini API
6. **结果融合**: 融合三层分数，生成最终结果

### 输出示例

```json
{
  "request_id": "20260127140530",
  "score": {
    "total": 68.5,
    "dimensions": {
      "vocabulary": 15.2,
      "sentence": 12.8,
      "personalization": 18.5,
      "logic": 14.0,
      "emotion": 8.0
    }
  },
  "multimodal": {
    "rule_layer_score": 65.0,
    "statistics_layer_score": 70.0,
    "semantic_layer_score": 72.0,
    "final_score": 68.5,
    "confidence": 0.78,
    "detection_mode": "multimodal"
  }
}
```

## API 配置参数

### 配置项说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | false | 是否启用 Gemini API |
| `api_key` | string | "" | API Key（建议通过环境变量设置） |
| `model` | string | "gemini-pro" | 使用的模型 |
| `temperature` | float | 0.3 | 生成温度（0-1，越低越确定） |
| `max_tokens` | int | 500 | 最大生成 token 数 |
| `timeout` | duration | 30s | API 调用超时时间 |

### 模型选择

目前支持的模型：

- **gemini-pro**: 标准模型，适合大多数场景
- **gemini-pro-vision**: 支持图像输入（未来版本）

### 温度参数

- **0.0-0.3**: 低温度，输出更确定、一致
- **0.4-0.7**: 中等温度，平衡创造性和一致性
- **0.8-1.0**: 高温度，输出更有创造性（不推荐用于检测）

**推荐**: 使用 0.3 以获得稳定的检测结果。

## 成本估算

### API 定价

Gemini API 定价（截至 2026-01）：

| 模型 | 输入 | 输出 | 说明 |
|------|------|------|------|
| gemini-pro | $0.00025/1K tokens | $0.0005/1K tokens | 标准定价 |
| gemini-pro (免费层) | 60 次/分钟 | 免费 | 有配额限制 |

### 单次检测成本

假设平均文本长度 1000 字：

- **输入 tokens**: ~1500 tokens（文本 + prompt）
- **输出 tokens**: ~200 tokens（分析结果）
- **单次成本**: ~$0.001（约 0.007 元人民币）

### 月度成本估算

| 检测量 | 月度成本 | 说明 |
|--------|---------|------|
| 1,000 次 | $1 | 小规模使用 |
| 10,000 次 | $10 | 中等规模 |
| 100,000 次 | $100 | 大规模使用 |

**注意**: 由于分层触发策略，实际 API 调用次数约为总检测次数的 20-30%。

## 错误处理

### 常见错误

#### 1. API Key 无效

**错误信息**: `Invalid API key`

**解决方法**:
- 检查 API Key 是否正确
- 确认 API Key 已启用
- 验证环境变量设置

```bash
# 检查环境变量
echo $GEMINI_API_KEY
```

#### 2. 配额超限

**错误信息**: `Quota exceeded`

**解决方法**:
- 等待配额重置（通常每分钟重置）
- 升级到付费计划
- 优化检测频率

#### 3. 网络超时

**错误信息**: `Request timeout`

**解决方法**:
- 检查网络连接
- 增加超时时间（配置文件中的 `timeout` 参数）
- 重试请求

```yaml
gemini:
  timeout: 60s  # 增加到60秒
```

#### 4. 服务不可用

**错误信息**: `Service unavailable` 或 `503 Error`

**解决方法**:
- 等待几分钟后重试
- 检查 Google API 服务状态
- 启用自动降级模式

### 自动重试机制

系统内置了自动重试机制:

```yaml
gemini:
  retry:
    max_attempts: 3      # 最大重试次数
    backoff: 1s          # 退避时间
    backoff_multiplier: 2 # 退避倍数
```

**重试策略**:
- 第1次失败: 等待 1 秒后重试
- 第2次失败: 等待 2 秒后重试
- 第3次失败: 等待 4 秒后重试
- 超过最大次数: 返回错误或降级

### 降级策略

当 Gemini API 不可用时,系统会自动降级:

```
完整多模态 → 规则+统计分析 → 仅规则检测
```

**降级触发条件**:
- API Key 未配置
- API 连续失败超过阈值
- 网络不可达
- 配额耗尽

**降级行为**:
- 自动禁用语义分析层
- 使用规则+统计分析继续检测
- 在结果中标注降级状态
- 记录降级原因到日志

## 最佳实践

### 1. API Key 安全管理

**推荐做法**:
```bash
# 使用环境变量
export GEMINI_API_KEY=your_api_key_here

# 添加到 .bashrc 或 .zshrc
echo 'export GEMINI_API_KEY=your_api_key_here' >> ~/.bashrc
```

**避免做法**:
- ❌ 不要在代码中硬编码 API Key
- ❌ 不要将 API Key 提交到版本控制
- ❌ 不要在配置文件中明文存储

### 2. 成本控制

**策略**:
- 使用分层触发策略,减少不必要的 API 调用
- 启用结果缓存,避免重复检测相同内容
- 设置每日调用限额

```yaml
gemini:
  cache:
    enabled: true
    ttl: 1h           # 缓存有效期
  rate_limit:
    max_per_day: 1000 # 每日最大调用次数
```

### 3. 性能优化

**建议**:
- 合理设置超时时间(推荐 30-60 秒)
- 使用低温度参数(0.3)保证结果一致性
- 限制 max_tokens 避免过长响应
- 启用并发检测(批量文件时)

```yaml
gemini:
  temperature: 0.3    # 低温度,结果更稳定
  max_tokens: 500     # 限制响应长度
  timeout: 30s        # 合理超时时间
```

### 4. 监控和日志

**监控指标**:
- API 调用成功率
- 平均响应时间
- 每日调用次数
- 成本统计

**日志记录**:
```bash
# 启用详细日志
aigc-check -f sample.txt -m -s -g --verbose > detection.log 2>&1
```

## 故障排查

### 问题诊断流程

1. **检查 API Key**
```bash
# 验证环境变量
echo $GEMINI_API_KEY

# 测试 API 连接
curl -H "Authorization: Bearer $GEMINI_API_KEY" \
  https://generativelanguage.googleapis.com/v1/models
```

2. **检查网络连接**
```bash
# 测试网络可达性
ping -c 3 generativelanguage.googleapis.com

# 检查代理设置
echo $HTTP_PROXY
echo $HTTPS_PROXY
```

3. **查看详细日志**
```bash
# 启用详细模式查看完整错误信息
aigc-check -f sample.txt -m -s -g --verbose
```

### 常见问题

**Q: 为什么有时候不调用 Gemini API?**

A: 系统使用分层触发策略。当规则检测或统计分析的置信度足够高时(>0.6),不会触发语义分析层,以节省成本。

**Q: 如何强制使用完整多模态检测?**

A: 目前系统会根据置信度自动决策。如需强制使用,可以调整配置文件中的置信度阈值:

```yaml
multimodal:
  confidence_thresholds:
    high: 0.95   # 提高阈值,更容易触发下一层
    medium: 0.80
    low: 0.50
```

**Q: API 调用失败后会怎样?**

A: 系统会自动重试(默认3次),如果仍然失败,会降级到规则+统计分析模式,确保检测能够继续进行。

**Q: 如何查看实际的 API 调用情况?**

A: 使用 `--verbose` 参数可以看到详细的分析过程,包括是否调用了 Gemini API:

```bash
aigc-check -f sample.txt -m -s -g --verbose
```

输出会显示:
- 规则层分数和置信度
- 是否触发统计分析
- 是否触发语义分析
- 最终融合结果

**Q: 缓存如何工作?**

A: 系统会缓存 Gemini API 的分析结果(默认1小时)。相同文本在缓存有效期内不会重复调用 API,直接返回缓存结果。

**Q: 如何估算月度成本?**

A: 根据实际使用情况:
- 统计每月检测次数
- 估算触发 API 的比例(通常 20-30%)
- 单次成本约 $0.001
- 月度成本 = 检测次数 × 触发比例 × 单次成本

示例: 10,000 次检测 × 25% × $0.001 = $2.5/月

## 参考资源

### 官方文档

- [Google AI Studio](https://makersuite.google.com/app/apikey) - 获取 API Key
- [Gemini API 文档](https://ai.google.dev/docs) - 官方 API 文档
- [Gemini API 定价](https://ai.google.dev/pricing) - 最新定价信息

### 相关文档

- [多模态检测架构](./multimodal-architecture.md) - 完整架构设计文档
- [README](../README.md) - 项目使用说明
- [配置文件示例](../configs/aigc-check.yaml) - 完整配置示例

### 技术支持

如遇到问题,请:
1. 查看本文档的故障排查部分
2. 检查 [GitHub Issues](https://github.com/leoobai/aigc-check/issues)
3. 提交新的 Issue 并附上详细日志

---

**文档版本**: v2.0
**更新日期**: 2026-01-27
**作者**: leoobai

