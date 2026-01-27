# AIGC-Check

AI生成内容检测工具 - 基于10个信号的智能检测系统

## 项目状态

🚧 **开发中** - MVP阶段

## 功能特性

- ✅ 10个AI生成内容检测信号
- ✅ 5维度评分系统
- ✅ 智能改进建议
- ✅ 多种输出格式（文本、JSON）
- ✅ 可配置的规则阈值

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/leoobai/aigc-check.git
cd aigc-check

# 构建
make build

# 安装到系统
make install
```

### 使用

```bash
# 检测文本文件
aigc-check detect -f sample.txt

# 使用JSON输出
aigc-check detect -f sample.txt -o json

# 使用自定义配置
aigc-check detect -f sample.txt -c custom-config.yaml
```

## 检测信号

1. **高频词汇** - 检测AI常用的关键词（crucial, pivotal等）
2. **句式开头** - 检测重复的句式开头（Additionally, Furthermore等）
3. **虚假范围** - 检测不连续的"from X to Y"表达
4. **引用异常** - 检测UTM参数和幽灵标记
5. **破折号密度** - 统计破折号使用频率
6. **Markdown残留** - 检测未清理的Markdown格式
7. **表情符号** - 检测过度工整的emoji使用
8. **知识截止** - 检测"截至我的知识更新"等短语
9. **协作式语气** - 检测"希望这能帮到你"等短语
10. **完美主义** - 检测缺乏第一人称和情感表达

## 评分维度

- **词汇多样性** (20分) - 评估词汇使用的丰富程度
- **句式复杂度** (15分) - 评估句式结构的多样性
- **个人化表达** (25分) - 评估个人风格和主观表达
- **逻辑连贯性** (20分) - 评估逻辑结构的自然性
- **情感真实度** (20分) - 评估情感表达的真实性

## 风险等级

- **极高风险** (0-40分) - 极可能为AI生成内容
- **高风险** (41-60分) - 很可能为AI生成内容
- **中等风险** (61-75分) - 可能包含AI生成内容
- **低风险** (76-100分) - 可能为人类编写

## 开发

```bash
# 运行测试
make test

# 查看测试覆盖率
make coverage

# 代码检查
make lint

# 清理构建产物
make clean
```

## 技术栈

- Go 1.21+
- gopkg.in/yaml.v3 - 配置解析

## 路线图

### Phase 1: MVP (当前)
- [x] 核心数据模型
- [ ] 规则引擎实现
- [ ] 评分系统
- [ ] CLI工具
- [ ] 基础测试

### Phase 2: 增强
- [ ] Gemini API集成
- [ ] 机器学习模型
- [ ] Web界面
- [ ] API服务

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License

## 作者

leoobai
