#!/bin/bash
set -e

# 服务启动脚本（生产部署：/data/media_report）
# 由 systemd 单元 media-report.service 调用，也可直接前台执行用于调试
# 关键点：cwd 必须切到 service/api，配置中的相对路径（etc/media-api.yaml、
# ./download_files、../../scripts/ulink、../../logs/sql.log 等）才能解析正确

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN="$ROOT_DIR/bin/media-api"
CONF="etc/media-api.yaml"

if [[ ! -x "$BIN" ]]; then
    echo "二进制不存在或不可执行: $BIN"
    echo "先在项目根执行: make build-api"
    exit 1
fi

cd "$ROOT_DIR/service/api"
exec "$BIN" -f "$CONF"
