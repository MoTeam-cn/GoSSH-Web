#!/bin/bash

# 获取当前版本号
get_current_version() {
  git tag -l "v*" | sort -V | tail -n1 | sed 's/^v//' || echo "0.9.0"
}

# 递增版本号
increment_version() {
  local version=$1
  local increment_type=$2

  IFS='.' read -ra parts <<< "$version"
  
  case $increment_type in
    "major")
      parts[0]=$((parts[0] + 1))
      parts[1]=0
      parts[2]=0
      ;;
    "minor")
      parts[1]=$((parts[1] + 1))
      parts[2]=0
      ;;
    "patch")
      parts[2]=$((parts[2] + 1))
      ;;
  esac

  echo "${parts[0]}.${parts[1]}.${parts[2]}"
}

# 生成版本日志
generate_changelog() {
  local new_version=$1
  local prev_tag=$2
  local current_date=$(date +"%Y-%m-%d")
  
  echo "# GoSSH-Web 版本 v${new_version} (${current_date})"
  echo ""
  
  if [ -z "$prev_tag" ]; then
    echo "## 初始版本"
    echo ""
    echo "* 初始化项目"
  else
    echo "## 更新内容"
    echo ""
    
    # 获取提交信息并格式化
    git log --pretty=format:"* %s" ${prev_tag}..HEAD | grep -v "Merge" | while read line; do
      if [[ "$line" == *"[feat]"* ]]; then
        echo "### ✨ 新功能"
        echo ""
        echo "${line/\[feat\] /}"
      elif [[ "$line" == *"[fix]"* ]]; then
        echo "### 🐛 修复"
        echo ""
        echo "${line/\[fix\] /}"
      elif [[ "$line" == *"[docs]"* ]]; then
        echo "### 📚 文档"
        echo ""
        echo "${line/\[docs\] /}"
      elif [[ "$line" == *"[chore]"* ]]; then
        echo "### 🔧 维护"
        echo ""
        echo "${line/\[chore\] /}"
      else
        echo "### 🔄 其他更新"
        echo ""
        echo "$line"
      fi
      echo ""
    done
  fi
}

# 主流程
main() {
  local current_version=$(get_current_version)
  local increment_type=${1:-"patch"}
  local new_version=$(increment_version "$current_version" "$increment_type")
  local prev_tag="v$current_version"
  
  # 检查是否有之前的标签
  if ! git tag -l | grep -q "$prev_tag"; then
    prev_tag=""
  fi
  
  # 生成变更日志
  generate_changelog "$new_version" "$prev_tag" > "CHANGELOG.md"
  
  echo "已生成变更日志 CHANGELOG.md 文件，新版本: v$new_version"
  echo "请检查变更日志内容，然后执行以下命令以提交并标记新版本:"
  echo "git add CHANGELOG.md"
  echo "git commit -m \"[build] Release v$new_version\""
  echo "git tag v$new_version"
  echo "git push origin main"
  echo "git push origin v$new_version"
}

# 执行主函数
main "$@" 