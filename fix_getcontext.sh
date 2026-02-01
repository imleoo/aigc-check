#!/bin/bash

# 批量修复所有规则文件中的 getContext 函数

files=(
  "internal/rules/signal01_high_freq_words.go"
  "internal/rules/signal04_citation_anomaly.go"
  "internal/rules/signal05_em_dash.go"
  "internal/rules/signal06_markdown.go"
  "internal/rules/signal07_emoji.go"
  "internal/rules/signal08_knowledge_cutoff.go"
  "internal/rules/signal09_collaborative.go"
)

for file in "${files[@]}"; do
  echo "Processing $file..."

  # 使用 sed 替换 getContext 函数
  sed -i.bak '/^func.*getContext/,/^}$/{
    /const contextSize = 50/a\
\
	// 将字节偏移量转换为 rune 索引
    /runes := \[\]rune(text)/a\
	runeOffset := len([]rune(text[:offset]))\
	runeLength := len([]rune(text[offset : offset+length]))
    s/start := offset - contextSize/start := runeOffset - contextSize/
    s/end := offset + length + contextSize/end := runeOffset + runeLength + contextSize/
  }' "$file"

  echo "Fixed $file"
done

echo "All files processed!"
