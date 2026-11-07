#!/usr/bin/env bash
# 依次运行所有已实现场景的自测（需要 Go 工具链）
set -euo pipefail
cd "$(dirname "$0")/.."

for d in advanced/*/; do
  [ -f "${d}main.go" ] || continue
  printf '\n========== %s ==========\n' "$d"
  go run "./${d%/}"
done

printf '\n全部场景自测通过\n'
