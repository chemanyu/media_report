#!/usr/bin/env bash
#
# media_report 部署脚本 —— 推到 172.16.3.34（从节点）
#
# 用法：
#   ./deploy.sh              # 完整流程：编译 → rsync → 远端安装/重启
#   ./deploy.sh build        # 仅交叉编译 linux 二进制
#   ./deploy.sh push         # 仅 rsync（不编译，不重启）
#   ./deploy.sh restart      # 仅在远端 systemctl restart
#   ./deploy.sh status       # 看远端服务状态
#   ./deploy.sh logs         # 远端 journalctl -f
#
# 环境变量覆盖：
#   REMOTE_USER=root   REMOTE_HOST=172.16.3.34   REMOTE_PATH=/data/media_report
#
# 重要：永远不会推送 service/api/etc/media-api.yaml
#       （远端 MySQL 密码、SyncFromProd / SyncDump 等配置与本机不同，请在 34 上手动维护）
#
set -euo pipefail

# ------ 配置 ------
REMOTE_USER="${REMOTE_USER:-root}"
REMOTE_HOST="${REMOTE_HOST:-172.16.3.34}"
REMOTE_PATH="${REMOTE_PATH:-/data/media_report}"

SERVICE_NAME="media-report"
LOCAL_ROOT="$(cd "$(dirname "$0")" && pwd)"
SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=10"

# ------ 颜色 ------
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'
info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERR]${NC}   $*" >&2; }

# ------ 步骤 ------

build_binary() {
    info "交叉编译 linux/amd64 二进制 → bin/media-api"
    (
        cd "$LOCAL_ROOT/service/api"
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
            go build -o "$LOCAL_ROOT/bin/media-api" media.go
    )
    file "$LOCAL_ROOT/bin/media-api" | grep -q 'ELF 64-bit.*x86-64' \
        || { error "编译产物不是 linux/amd64 ELF"; exit 1; }
    info "编译完成: $(ls -lh "$LOCAL_ROOT/bin/media-api" | awk '{print $5,$9}')"
}

push_files() {
    info "rsync → $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"
    [[ ! -x "$LOCAL_ROOT/bin/media-api" ]] && { error "bin/media-api 不存在，先跑 build"; exit 1; }

    local excludes=(
        --exclude='.git/'
        --exclude='.idea/'
        --exclude='.vscode/'
        --exclude='.claude/'
        --exclude='logs/'
        --exclude='uploads/'
        --exclude='download_files/'
        --exclude='*.log'
        --exclude='coverage.*'
        --exclude='.DS_Store'
        --exclude='dist/'
        # 永不推送配置文件：远端 MySQL 凭据 / SyncFromProd / SyncDump 与本机不同
        --exclude='service/api/etc/media-api.yaml'
    )

    warn "etc/media-api.yaml 不会被同步（受 deploy.sh 硬编码排除）；如需修改远端配置请 ssh 上去手动改"

    rsync -avz --delete-after \
        "${excludes[@]}" \
        -e "ssh $SSH_OPTS" \
        "$LOCAL_ROOT/" \
        "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH/"
}

remote_setup() {
    info "远端安装 / 重启 systemd 服务"
    ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" bash <<EOF
set -euo pipefail
cd "$REMOTE_PATH"

# 1. 给二进制 + 启动脚本加执行权限
chmod +x bin/media-api start-api.sh

# 2. 准备日志目录
mkdir -p logs download_files uploads

# 3. 安装 / 更新 systemd unit
UNIT_SRC="$REMOTE_PATH/deploy/${SERVICE_NAME}.service"
UNIT_DST="/etc/systemd/system/${SERVICE_NAME}.service"
if [[ ! -f "\$UNIT_SRC" ]]; then
    echo "ERROR: \$UNIT_SRC 不存在"
    exit 1
fi
if ! cmp -s "\$UNIT_SRC" "\$UNIT_DST" 2>/dev/null; then
    echo "[INFO] systemd unit 有变化，更新中..."
    cp "\$UNIT_SRC" "\$UNIT_DST"
    systemctl daemon-reload
else
    echo "[INFO] systemd unit 无变化"
fi

# 4. 启用 + 重启
systemctl enable ${SERVICE_NAME} >/dev/null 2>&1 || true
if systemctl is-active --quiet ${SERVICE_NAME}; then
    echo "[INFO] 服务在跑，执行 restart"
    systemctl restart ${SERVICE_NAME}
else
    echo "[INFO] 服务未启动，执行 start"
    systemctl start ${SERVICE_NAME}
fi

sleep 2
systemctl --no-pager status ${SERVICE_NAME} | head -20
EOF
}

remote_status() {
    ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" \
        "systemctl --no-pager status ${SERVICE_NAME}"
}

remote_logs() {
    ssh $SSH_OPTS -t "$REMOTE_USER@$REMOTE_HOST" \
        "journalctl -u ${SERVICE_NAME} -f"
}

remote_restart() {
    ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" \
        "systemctl restart ${SERVICE_NAME} && sleep 1 && systemctl --no-pager status ${SERVICE_NAME} | head -10"
}

# ------ 入口 ------
cmd="${1:-deploy}"
case "$cmd" in
    build)   build_binary ;;
    push)    push_files ;;
    restart) remote_restart ;;
    status)  remote_status ;;
    logs)    remote_logs ;;
    deploy|"")
        build_binary
        push_files
        remote_setup
        info "部署完成: http://$REMOTE_HOST:8888"
        info "查看日志: ./deploy.sh logs"
        ;;
    *)
        echo "未知命令: $cmd"
        echo "用法: $0 [build|push|restart|status|logs|deploy]"
        exit 1
        ;;
esac
