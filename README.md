# Bootcs Schema

Bootcs 课程和关卡配置的 JSON Schema 定义。

## 📋 Schema 文件

- **`course.schema.json`** - 课程配置 (`course.yml`)
- **`stage.schema.json`** - 关卡配置 (`stage.yml`)

## 🚀 使用方法

### 在 YAML 文件中引用 Schema

**course.yml**:

```yaml
# yaml-language-server: $schema=https://bootcs.dev/schemas/course.schema.json

slug: my-course
name: "我的课程"
summary: "课程简介"
```

**stage.yml**:

```yaml
# yaml-language-server: $schema=https://bootcs.dev/schemas/stage.schema.json

slug: hello
name: "Hello"
summary: "关卡简介"
description: "README.md"
files:
  required: ["hello.c"]
  allowed: ["*.c", "*.h"]
```

### VS Code 支持

安装 [YAML 扩展](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)：

```bash
code --install-extension redhat.vscode-yaml
```

获得实时验证、自动补全和文档提示。
