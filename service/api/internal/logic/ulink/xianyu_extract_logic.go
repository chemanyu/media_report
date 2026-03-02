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
	"media_report/service/api/internal/types"
)

// XianyuExtractLogic 闲鱼 deeplink 提取（单链 + 批量）
type XianyuExtractLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXianyuExtractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XianyuExtractLogic {
	return &XianyuExtractLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// xianyuScriptResult Python 脚本单链输出
type xianyuScriptResult struct {
	Deeplink string `json:"deeplink"`
	Error    string `json:"error"`
}

// ExtractSingle 单链提取
func (l *XianyuExtractLogic) ExtractSingle(req *types.XianyuExtractReq) (*types.XianyuExtractResp, error) {
	cfg := l.svcCtx.Config.Ulink
	platform := req.Platform
	if platform == "" {
		platform = "android"
	}

	scriptPath := filepath.Join(cfg.ScriptDir, "xianyu_deeplink.py")
	args := []string{
		"--mode", "single",
		"--url", req.ShortUrl,
		"--platform", platform,
	}
	if cfg.ChromeDriverPath != "" {
		args = append(args, "--driver", cfg.ChromeDriverPath)
	}

	out, err := runScript(cfg.PythonPath, scriptPath, args...)
	if err != nil {
		l.Errorf("调用 xianyu_deeplink.py 失败: %v, output: %s", err, string(out))
		return &types.XianyuExtractResp{Code: 500, Message: "调用 Python 脚本失败: " + err.Error()}, nil
	}

	var result xianyuScriptResult
	if err := json.Unmarshal(lastJSONLine(out), &result); err != nil {
		l.Errorf("解析脚本输出失败: %v, output: %s", err, string(out))
		return &types.XianyuExtractResp{Code: 500, Message: "解析脚本输出失败"}, nil
	}

	if result.Error != "" {
		return &types.XianyuExtractResp{Code: 400, Message: result.Error}, nil
	}

	return &types.XianyuExtractResp{
		Code:     200,
		Message:  "success",
		Deeplink: result.Deeplink,
	}, nil
}

// ExtractBatch 批量提取
func (l *XianyuExtractLogic) ExtractBatch(w http.ResponseWriter, r *http.Request) {
	cfg := l.svcCtx.Config.Ulink
	platform := r.FormValue("platform")
	if platform == "" {
		platform = "android"
	}

	inputPath, cleanup, err := saveUploadedFile(r, "link_file", cfg.TempDir, "*.txt")
	if err != nil {
		writeJSONError(w, 400, "读取上传文件失败: "+err.Error())
		return
	}
	defer cleanup()

	outputPath := filepath.Join(cfg.TempDir, fmt.Sprintf("xianyu_batch_%d.xlsx", time.Now().UnixNano()))
	defer os.Remove(outputPath)

	scriptPath := filepath.Join(cfg.ScriptDir, "xianyu_deeplink.py")
	args := []string{
		"--mode", "batch",
		"--input", inputPath,
		"--output", outputPath,
		"--platform", platform,
	}
	if cfg.ChromeDriverPath != "" {
		args = append(args, "--driver", cfg.ChromeDriverPath)
	}

	out, err := runScript(cfg.PythonPath, scriptPath, args...)
	if err != nil {
		l.Errorf("闲鱼批量提取脚本执行失败: %v, output: %s", err, string(out))
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

	serveExcelFile(w, outputPath, "xianyu_deeplink_results.xlsx")
}
