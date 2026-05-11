# 部署到 172.16.2.34:/data/media_report

## 1. 编译二进制（在开发机或服务器都行）

```bash
cd /data/media_report
make build-api
# 产物: /data/media_report/bin/media-api
```

如果在 Mac/Windows 编译给 Linux 用：

```bash
cd service/api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../bin/media-api media.go
```

## 2. 安装 systemd 单元

```bash
sudo cp /data/media_report/deploy/media-report.service /etc/systemd/system/media-report.service
sudo systemctl daemon-reload
sudo systemctl enable media-report   # 开机自启
sudo systemctl start media-report
```

## 3. 常用操作

```bash
sudo systemctl status media-report      # 查看状态
sudo systemctl restart media-report     # 重启（更新二进制后）
sudo systemctl stop media-report
sudo journalctl -u media-report -f      # 跟随 systemd 日志
tail -f /data/media_report/logs/media-api.stdout.log
tail -f /data/media_report/logs/media-api.stderr.log
tail -f /data/media_report/logs/sql.log
```

## 4. 调试（不通过 systemd）

```bash
cd /data/media_report
./start-api.sh
```

`start-api.sh` 会切到 `service/api/` 后用同一个二进制前台运行，便于排错。

## 5. 升级流程

```bash
cd /data/media_report
git pull
make build-api
sudo systemctl restart media-report
```

## 6. 常见坑

- **必须以 `service/api/` 作为 cwd**：`media-api.yaml` 里 `etc/media-api.yaml`、`./download_files`、`../../scripts/ulink`、`../../logs/sql.log` 全是相对路径，cwd 错了会读不到配置/写不到日志。
- **从节点（172.16.2.34）必须把 yaml 中 `SyncFromProd.Enabled` 设为 `true`、`SyncDump.Enabled` 设为 `false`**，并填好 `BaseURL` 和 `Token`，否则配置表不会自动从主生产同步。
- **Chrome 扩展 / Python 转链脚本** 在服务器侧不需要也用不到（无 GUI / 无 ChromeDriver），相关功能只在主生产或本地有 Chrome 的机器使用。
