package download

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/svc"
)

type DownloadReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadReportLogic {
	return &DownloadReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadReportLogic) DownloadReport(w http.ResponseWriter, r *http.Request) error {
	// 从URL路径中获取文件名
	filename := r.URL.Query().Get(":filename")
	if filename == "" {
		// 尝试从路径参数中获取
		pathSegments := filepath.Base(r.URL.Path)
		filename = pathSegments
	}

	// 防止路径遍历攻击
	filename = filepath.Base(filename)

	// 获取文件存储路径
	savePath := l.svcCtx.Config.FileServer.Path

	// 构建完整文件路径
	filePath := filepath.Join(savePath, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return nil
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		logx.Errorf("打开文件失败: %v", err)
		http.Error(w, "打开文件失败", http.StatusInternalServerError)
		return nil
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		logx.Errorf("获取文件信息失败: %v", err)
		http.Error(w, "获取文件信息失败", http.StatusInternalServerError)
		return nil
	}

	// 设置响应头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 将文件内容写入响应
	http.ServeContent(w, r, filename, fileInfo.ModTime(), file)

	logx.Infof("文件下载成功: %s", filename)
	return nil
}
