# 贡献指南

感谢您对 WebTTY 项目感兴趣！我们欢迎并感谢任何形式的贡献。

## 如何贡献

### 报告问题

1. 使用 GitHub Issues 提交问题
2. 清晰描述问题，包括：
   - 问题的具体表现
   - 复现步骤
   - 期望的行为
   - 实际的行为
   - 环境信息（操作系统、浏览器等）
3. 如果可能，提供截图或录屏

### 提交代码

1. Fork 项目
2. 创建特性分支
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. 提交更改
   ```bash
   git commit -m "Add some feature"
   ```
4. 推送到分支
   ```bash
   git push origin feature/your-feature-name
   ```
5. 创建 Pull Request

### 代码规范

- 遵循 Go 标准代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 编写单元测试
- 确保所有测试通过

### 提交信息规范

提交信息应该清晰描述更改内容，建议使用以下格式：

- feat: 新功能
- fix: 修复问题
- docs: 文档更新
- style: 代码格式（不影响代码运行的变动）
- refactor: 重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

示例：
```
feat: 添加 SSH 密钥认证支持
```

### 文档贡献

- 改进现有文档
- 添加使用示例
- 修正文档错误
- 翻译文档

## 开发设置

1. 安装依赖
   ```bash
   go mod tidy
   ```

2. 运行测试
   ```bash
   go test ./...
   ```

3. 启动开发服务器
   ```bash
   go run main.go
   ```

## 行为准则

- 尊重所有贡献者
- 保持专业和友善的交流
- 关注技术讨论
- 接受建设性的批评

## 获取帮助

- 查看项目文档
- 提交 Issue
- 通过 Pull Request 讨论

再次感谢您的贡献！ 