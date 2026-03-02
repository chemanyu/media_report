package ulink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/svc"
)

// ActivityBatchLogic 淘宝客活动短链批量获取
type ActivityBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivityBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivityBatchLogic {
	return &ActivityBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetActivityBatch 活动短链批量获取（文件上传 → Excel 下载）
func (l *ActivityBatchLogic) GetActivityBatch(w http.ResponseWriter, r *http.Request) {
	cfg := l.svcCtx.Config.Ulink

	materialId := r.FormValue("activity_material_id")
	if materialId == "" {
		materialId = cfg.DefaultMaterialId
	}

	// 读取上传的 sub_pid 列表文件
	inputPath, cleanup, err := saveUploadedFile(r, "link_file", cfg.TempDir, "*.txt")
	if err != nil {
		writeJSONError(w, 400, "读取上传文件失败: "+err.Error())
		return
	}
	defer cleanup()

	outputPath := filepath.Join(cfg.TempDir, fmt.Sprintf("activity_batch_%d.xlsx", time.Now().UnixNano()))
	defer os.Remove(outputPath)

	scriptPath := filepath.Join(cfg.ScriptDir, "activity_batch.py")
	args := []string{
		"--input", inputPath,
		"--output", outputPath,
		"--material_id", materialId,
		"--app_key", cfg.TaobaoAppKey,
		"--app_secret", cfg.TaobaoAppSecret,
	}
	if cfg.ChromeDriverPath != "" {
		args = append(args, "--driver", cfg.ChromeDriverPath)
	}

	out, err := runScript(cfg.PythonPath, scriptPath, args...)
	if err != nil {
		l.Errorf("活动短链批量脚本执行失败: %v, output: %s", err, string(out))
		writeJSONError(w, 500, "脚本执行失败: "+err.Error())
		return
	}

	var result batchScriptResult
	if err := json.Unmarshal(lastJSONLine(out), &result); err != nil {
		writeJSONError(w, 500, "解析脚本输出失败")
		return
	}

	if result.Error != "" {
		writeJSONError(w, 400, result.Error)
		return
	}

	serveExcelFile(w, outputPath, "activity_shortlinks.xlsx")
}
