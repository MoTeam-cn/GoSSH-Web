# 版本发布脚本

# 检查是否提供了提交信息
param(
    [string]$Message = "维护性更新"
)

# 添加所有文件
git add .

# 提交更改
git commit -m "[build] $Message"

# 推送到主分支
git push origin main

Write-Host "已触发构建流程，版本将自动更新。请等待GitHub Actions完成处理。"
Write-Host "构建进度可以在此查看: https://github.com/MoTeam-cn/GoSSH-Web/actions" 