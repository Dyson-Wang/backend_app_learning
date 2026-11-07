#!/usr/bin/env bash
# 一键验证五门语言的 basics 24 课与 advanced 场景。
# 原则：能编译就编译，能运行就运行；缺工具链明确跳过，不伪装通过。
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
PASS=0
FAIL=0
SKIP=0

run() {
  printf '\n>>> %s\n' "$1"
  shift
  "$@"
  PASS=$((PASS + 1))
}

# 带超时运行（脚本类示例必须自行结束，卡住即失败）
run_limited() {
  local label="$1"
  shift
  printf '\n>>> %s\n' "$label"
  if perl -e 'alarm 8; exec @ARGV' "$@" >/tmp/brm-run.log 2>&1; then
    PASS=$((PASS + 1))
  else
    FAIL=$((FAIL + 1))
    printf '  FAILED（退出码 %s）：\n' "$?"
    sed 's/^/    /' /tmp/brm-run.log | tail -8
  fi
}

skip() {
  printf '\n[skip] %s\n' "$1"
  SKIP=$((SKIP + 1))
}

# ---------- Go ----------
if command -v go >/dev/null 2>&1; then
  run "Go: go build ./..." bash -c 'cd "$1/01-go" && go build ./...' _ "$ROOT"
  run "Go: go vet ./..." bash -c 'cd "$1/01-go" && go vet ./...' _ "$ROOT"
  for d in "$ROOT"/01-go/basics/*/; do
    [ -f "${d}main.go" ] || continue
    name="$(basename "$d")"
    [ "$name" = "10-first-http-service-with-db" ] && continue
    run_limited "Go: go run $name" bash -c 'cd "$1" && go run .' _ "$d"
  done
  if [ -x "$ROOT/01-go/advanced/run-selftests.sh" ]; then
    run "Go: advanced 场景自测" bash "$ROOT/01-go/advanced/run-selftests.sh"
  fi
else
  skip "Go 未安装（brew install go），跳过 01-go 编译与运行验证"
fi

# ---------- TypeScript / Node.js ----------
if command -v node >/dev/null 2>&1; then
  for f in "$ROOT"/02-typescript-node/basics/*/index.mjs; do
    [ -f "$f" ] || continue
    run "Node: node --check $(basename "$(dirname "$f")")" node --check "$f"
  done
  for f in "$ROOT"/02-typescript-node/basics/*/index.mjs; do
    [ -f "$f" ] || continue
    name="$(basename "$(dirname "$f")")"
    [ "$name" = "10-first-http-service-with-db" ] && continue
    run_limited "Node: 运行 $name" node "$f"
  done
  for f in "$ROOT"/02-typescript-node/advanced/*/index.mjs; do
    [ -f "$f" ] || continue
    run "Node: node --check advanced/$(basename "$(dirname "$f")")" node --check "$f"
    run_limited "Node: 运行 advanced/$(basename "$(dirname "$f")")" node "$f"
  done
else
  skip "Node 未安装"
fi

# ---------- Python ----------
if command -v python3 >/dev/null 2>&1; then
  for f in "$ROOT"/03-python/basics/*/main.py; do
    [ -f "$f" ] || continue
    run "Python: py_compile $(basename "$(dirname "$f")")" python3 -m py_compile "$f"
  done
  for f in "$ROOT"/03-python/basics/*/main.py; do
    [ -f "$f" ] || continue
    name="$(basename "$(dirname "$f")")"
    [ "$name" = "10-first-http-service-with-db" ] && continue
    run_limited "Python: 运行 $name" python3 -B "$f"
  done
  for f in "$ROOT"/03-python/advanced/*/main.py; do
    [ -f "$f" ] || continue
    run "Python: py_compile advanced/$(basename "$(dirname "$f")")" python3 -m py_compile "$f"
    run_limited "Python: 运行 advanced/$(basename "$(dirname "$f")")" python3 -B "$f"
  done
  find "$ROOT" -name __pycache__ -type d -exec rm -rf {} + 2>/dev/null
else
  skip "Python 未安装"
fi

# ---------- Java ----------
if command -v javac >/dev/null 2>&1 && javac -version >/dev/null 2>&1; then
  out="/tmp/brm-java-out"
  rm -rf "$out" && mkdir -p "$out"
  for d in "$ROOT"/04-java/basics/*/; do
    [ -f "${d}Main.java" ] || continue
    run "Java: javac $(basename "$d")" javac -d "$out" "${d}Main.java"
  done
  for d in "$ROOT"/04-java/advanced/*/; do
    [ -f "${d}Main.java" ] || continue
    run "Java: javac advanced/$(basename "$d")" javac -d "$out" "${d}Main.java"
  done
else
  skip "javac 不可用（macOS 自带 java 可能只是占位 stub），跳过 04-java 编译验证"
fi

# ---------- Rust ----------
if command -v cargo >/dev/null 2>&1; then
  for d in "$ROOT"/05-rust/basics/*/; do
    [ -f "${d}Cargo.toml" ] || continue
    name="$(basename "$d")"
    run "Rust: cargo check $name" bash -c 'cd "$1" && cargo check --quiet' _ "$d"
    [ "$name" = "10-first-http-service-with-db" ] && continue
    run_limited "Rust: cargo run $name" bash -c 'cd "$1" && cargo run --quiet' _ "$d"
  done
  for d in "$ROOT"/05-rust/advanced/*/; do
    [ -f "${d}Cargo.toml" ] || continue
    name="$(basename "$d")"
    run "Rust: cargo check advanced/$name" bash -c 'cd "$1" && cargo check --quiet' _ "$d"
    run_limited "Rust: cargo run advanced/$name" bash -c 'cd "$1" && cargo run --quiet' _ "$d"
  done
else
  skip "cargo 未安装，跳过 05-rust 编译与运行验证"
fi

printf '\n完成：%s 项通过，%s 项失败，%s 项跳过\n' "$PASS" "$FAIL" "$SKIP"
[ "$FAIL" -eq 0 ] || exit 1
