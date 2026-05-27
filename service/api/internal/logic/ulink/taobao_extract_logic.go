package ulink

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// TaobaoExtractLogic 淘宝 deeplink 提取（单链 + 批量）
type TaobaoExtractLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaobaoExtractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaobaoExtractLogic {
	return &TaobaoExtractLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// taobaoScriptResult Python 脚本单链输出的 JSON 结构
type taobaoScriptResult struct {
	Deeplink string `json:"deeplink"`
	H5Dp     string `json:"h5_dp"`
	Error    string `json:"error"`
}

// batchScriptResult Python 脚本批量输出的 JSON 结构
type batchScriptResult struct {
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

// ExtractSingle 单链提取
func (l *TaobaoExtractLogic) ExtractSingle(req *types.TaobaoExtractReq) (*types.TaobaoExtractResp, error) {
	cfg := l.svcCtx.Config.Ulink
	platform := req.Platform
	if platform == "" {
		platform = "ios"
	}

	scriptPath := filepath.Join(cfg.ScriptDir, "taobao_deeplink.py")
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
		l.Errorf("调用 taobao_deeplink.py 失败: %v, output: %s", err, string(out))
		return &types.TaobaoExtractResp{Code: 500, Message: "调用 Python 脚本失败: " + err.Error()}, nil
	}

	var result taobaoScriptResult
	if err := json.Unmarshal(lastJSONLine(out), &result); err != nil {
		l.Errorf("解析脚本输出失败: %v, output: %s", err, string(out))
		return &types.TaobaoExtractResp{Code: 500, Message: "解析脚本输出失败"}, nil
	}

	if result.Error != "" {
		return &types.TaobaoExtractResp{Code: 400, Message: result.Error}, nil
	}

	return &types.TaobaoExtractResp{
		Code:     200,
		Message:  "success",
		Deeplink: result.Deeplink,
		H5Dp:     result.H5Dp,
	}, nil
}

// ExtractBatch 批量提取（文件上传 → Excel 下载）
func (l *TaobaoExtractLogic) ExtractBatch(w http.ResponseWriter, r *http.Request) {
	cfg := l.svcCtx.Config.Ulink
	platform := r.FormValue("platform")
	if platform == "" {
		platform = "ios"
	}

	// 读取上传文件
	inputPath, cleanup, err := saveUploadedFile(r, "link_file", cfg.TempDir, "*.txt")
	if err != nil {
		writeJSONError(w, 400, "读取上传文件失败: "+err.Error())
		return
	}
	defer cleanup()

	// 生成输出文件路径
	outputPath := filepath.Join(cfg.TempDir, fmt.Sprintf("taobao_batch_%d.xlsx", time.Now().UnixNano()))
	defer os.Remove(outputPath)

	scriptPath := filepath.Join(cfg.ScriptDir, "taobao_deeplink.py")
	args := []string{
		"--mode", "batch",
		"--input", inputPath,
		"--output", outputPath,
		"--platform", platform,
	}
	if cfg.ChromeDriverPath != "" {
		args = append(args, "--driver", cfg.ChromeDriverPath)
	}

	out, err := runScriptWithContext(r.Context(), cfg.PythonPath, scriptPath, args...)
	if err != nil {
		l.Errorf("批量提取脚本执行失败: %v, output: %s", err, string(out))
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

	// 直接将 Excel 输出到 response
	serveExcelFile(w, outputPath, "taobao_deeplink_results.xlsx")
}

// applyPythonEnv 给 Python 子进程设置 UTF-8 环境，避免 Windows GBK 编码崩溃，
// 并把工作目录设到脚本所在目录，保证相对 import 正常。
func applyPythonEnv(cmd *exec.Cmd, scriptPath string) {
	cmd.Dir = filepath.Dir(scriptPath)
	cmd.Env = append(os.Environ(),
		"PYTHONIOENCODING=utf-8",
		"PYTHONUTF8=1",
	)
}

// runScript 调用 Python 脚本，绑定 context（超时/取消时自动终止子进程）
// 返回 stdout；失败时错误信息里附带 stderr，方便 Windows 调试。
func runScript(pythonPath, scriptPath string, args ...string) ([]byte, error) {
	allArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(pythonPath, allArgs...)
	applyPythonEnv(cmd, scriptPath)
	return runWithStderr(cmd)
}

// runScriptWithContext 调用 Python 脚本并绑定 context，超时后自动 kill 子进程
func runScriptWithContext(ctx context.Context, pythonPath, scriptPath string, args ...string) ([]byte, error) {
	allArgs := append([]string{scriptPath}, args...)
	cmd := exec.CommandContext(ctx, pythonPath, allArgs...)
	applyPythonEnv(cmd, scriptPath)
	return runWithStderr(cmd)
}

func runWithStderr(cmd *exec.Cmd) ([]byte, error) {
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("%w; stderr: %s", err, stderr.String())
	}
	return out, nil
}

// lastJSONLine 从输出中取最后一个 JSON 行（Python 脚本可能会打印日志，最后一行是 JSON）
func lastJSONLine(output []byte) []byte {
	lines := splitLines(output)
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if len(line) > 0 && line[0] == '{' {
			return line
		}
	}
	return output
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

// saveUploadedFile 将上传的文件保存到临时路径，返回路径和清理函数
func saveUploadedFile(r *http.Request, fieldName, tempDir, pattern string) (string, func(), error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", nil, err
	}
	file, _, err := r.FormFile(fieldName)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	tmpFile, err := os.CreateTemp(tempDir, pattern)
	if err != nil {
		return "", nil, err
	}
	tmpPath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", nil, err
	}
	tmpFile.Close()

	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

// serveExcelFile 将 Excel 文件写入 HTTP response
func serveExcelFile(w http.ResponseWriter, filePath, downloadName string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		writeJSONError(w, 500, "打开结果文件失败")
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// writeJSONError 写入 JSON 错误响应
func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	resp, _ := json.Marshal(map[string]interface{}{"code": code, "message": msg})
	w.Write(resp)
}
