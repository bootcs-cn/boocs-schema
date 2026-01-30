# 更新日志

## [0.3.0] - 2026-01-30

### 变更

- **BREAKING**: 将 stage schema 中的 `files` 字段重命名为 `files_config`
- **BREAKING**: 将 `evaluation.timeout` 移至顶层 `timeout` 字段
- **BREAKING**: 从 stage schema 中移除 `language` 字段（现在在 submissions 表中）
- 简化 schema 以匹配数据库迁移 v2.1
- 所有 30 个 bcs100x stages 通过验证

## [0.2.0] - 2026-01-29

### 新增

- 课程配置 schema (`course.schema.json`)
- 关卡配置 schema (`stage.schema.json`)
- 支持 JSON Schema Draft-07
- 中文描述和文档
