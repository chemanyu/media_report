package ulink

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/svc"
)

// adzoneGroup 单个广告位ID及其对应的商品列表
type adzoneGroup struct {
	rawId    string // Excel 原始行内容（用于写入导出文件）
	adzoneId string // 分割后的真实 adzone_id（用于 API 调用）
	items    []map[string]interface{}
}

// CpsGoodsLogic CPS高佣商品库批量导出
type CpsGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCpsGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CpsGoodsLogic {
	return &CpsGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCpsGoods 读取上传 Excel 中的 adzone_id 列表，分页拉取商品并输出分组 Excel
func (l *CpsGoodsLogic) GetCpsGoods(w http.ResponseWriter, r *http.Request) {
	cfg := l.svcCtx.Config.Ulink

	materialId := r.FormValue("material_id")
	if materialId == "" {
		writeJSONError(w, 400, "material_id 不能为空")
		return
	}

	// 读取上传的 Excel 文件
	inputPath, cleanup, err := saveUploadedFile(r, "link_file", cfg.TempDir, "*.xlsx")
	if err != nil {
		writeJSONError(w, 400, "读取上传文件失败: "+err.Error())
		return
	}
	defer cleanup()

	// 解析 A 列 adzone_id 列表
	adzoneIds, err := readColumnA(inputPath)
	if err != nil {
		writeJSONError(w, 400, "解析 Excel 失败: "+err.Error())
		return
	}
	if len(adzoneIds) == 0 {
		writeJSONError(w, 400, "Excel A 列未找到有效的 adzone_id")
		return
	}

	// 对每个 adzone_id 拉取全量商品
	groups := make([]adzoneGroup, 0, len(adzoneIds))

	for _, rawId := range adzoneIds {
		// 格式如 mm_874030133_3340450211_116227700211，取最后一段作为真实 adzone_id
		parts := strings.Split(rawId, "_")
		adzoneId := parts[len(parts)-1]

		var allItems []map[string]interface{}
		pageNo := 1
		pageSize := 100
		for {
			items, err := fetchCpsMaterialPage(cfg.TaobaoAppKey, cfg.TaobaoAppSecret, materialId, adzoneId, pageNo, pageSize)
			if err != nil {
				l.Errorf("adzone_id=%s 第%d页拉取失败: %v", adzoneId, pageNo, err)
				break
			}
			if len(items) == 0 {
				break
			}
			allItems = append(allItems, items...)
			if len(items) < pageSize {
				break
			}
			pageNo++
		}
		groups = append(groups, adzoneGroup{rawId: rawId, adzoneId: adzoneId, items: allItems})
	}

	// 生成分组 Excel
	outputPath := filepath.Join(cfg.TempDir, fmt.Sprintf("cps_goods_%d.xlsx", time.Now().UnixNano()))
	defer os.Remove(outputPath)

	if err := buildCpsGoodsGroupedExcel(outputPath, groups); err != nil {
		l.Errorf("生成CPS商品库Excel失败: %v", err)
		writeJSONError(w, 500, "生成Excel失败: "+err.Error())
		return
	}

	serveExcelFile(w, outputPath, fmt.Sprintf("cps_goods_%s.xlsx", materialId))
}

// readColumnA 从 Excel 文件读取第一个 Sheet 的 A 列，返回非空行的值列表
func readColumnA(filePath string) ([]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 文件没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, row := range rows {
		if len(row) > 0 {
			val := strings.TrimSpace(row[0])
			if val != "" {
				ids = append(ids, val)
			}
		}
	}
	return ids, nil
}

// fetchCpsMaterialPage 调用 taobao.tbk.dg.material.recommend 接口拉取一页商品
func fetchCpsMaterialPage(appKey, appSecret, materialId, adzoneId string, pageNo, pageSize int) ([]map[string]interface{}, error) {
	params := map[string]string{
		"method":      "taobao.tbk.dg.material.recommend",
		"app_key":     appKey,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"format":      "json",
		"v":           "2.0",
		"sign_method": "md5",
		"partner_id":  "top-apitools",
		"material_id": materialId,
		"adzone_id":   adzoneId,
		"page_no":     fmt.Sprintf("%d", pageNo),
		"page_size":   fmt.Sprintf("%d", pageSize),
	}
	params["sign"] = signTaobaoAPI(params, appSecret)

	queryValues := url.Values{}
	for k, v := range params {
		queryValues.Set(k, v)
	}

	resp, err := taobaoHTTPClient.Get(taobaoAPIURL + "?" + queryValues.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if errResp, hasErr := result["error_response"]; hasErr {
		errMap, _ := errResp.(map[string]interface{})
		msg := ""
		if errMap != nil {
			msg = fmt.Sprintf("%v", errMap["msg"])
		}
		return nil, fmt.Errorf("淘宝API错误: %s", msg)
	}

	respData, ok := result["tbk_dg_material_recommend_response"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	resultList, ok := respData["result_list"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	rawItems, ok := resultList["map_data"].([]interface{})
	if !ok {
		return nil, nil
	}

	items := make([]map[string]interface{}, 0, len(rawItems))
	for _, raw := range rawItems {
		if item, ok := raw.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items, nil
}

// cpsGoodsHeaders Excel 表头（与 buildCpsGoodsRow 中的取值顺序严格对应）
var cpsGoodsHeaders = []string{
	// --- 标识 ---
	"adzone_id",
	// --- 商品基本信息 ---
	"商品ID",
	"标题",
	"短标题",
	"子标题",
	"品牌",
	"叶子分类",
	"一级分类",
	"店铺类型", // 0=淘宝 1=天猫
	"店铺名称",
	"月销量",
	"年销量",
	"运费",
	// --- 价格 & 优惠券 ---
	"原价",
	"折扣价",
	"最终成交价",
	"券类型",
	"优惠券面额",
	"优惠券使用条件",
	"券开始时间",
	"券结束时间",
	"凑单价",
	"凑单描述",
	// --- 推广信息 ---
	"佣金率(%)",
	// --- 链接 ---
	"主图",
	"推广链接",
	"优惠券推广链接",
}

// buildCpsGoodsGroupedExcel 生成分组商品 Excel：每个 adzone_id 行后紧跟其商品行
func buildCpsGoodsGroupedExcel(outputPath string, groups []adzoneGroup) error {
	f := excelize.NewFile()
	sheet := "CPS高佣商品库"
	f.SetSheetName("Sheet1", sheet)

	// 写表头（第 1 行）
	for col, h := range cpsGoodsHeaders {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	currentRow := 2
	for _, g := range groups {
		// adzone_id 分组行（仅 A 列），写原始值
		cell, _ := excelize.CoordinatesToCellName(1, currentRow)
		f.SetCellValue(sheet, cell, g.rawId)
		currentRow++

		// 商品数据行（A 列留空，B 列起）
		for _, item := range g.items {
			values := buildCpsGoodsRow(g.adzoneId, item)
			for col, v := range values {
				cell, _ := excelize.CoordinatesToCellName(col+1, currentRow)
				if v != nil && v != "" {
					f.SetCellValue(sheet, cell, v)
				}
			}
			currentRow++
		}
	}

	return f.SaveAs(outputPath)
}

// buildCpsGoodsRow 将单条商品 API 响应展开为与 cpsGoodsHeaders 对应的值切片
// A 列（adzone_id）留空，由分组行已写出
func buildCpsGoodsRow(_ string, item map[string]interface{}) []interface{} {
	basic, _ := item["item_basic_info"].(map[string]interface{})
	price, _ := item["price_promotion_info"].(map[string]interface{})
	publish, _ := item["publish_info"].(map[string]interface{})

	// 店铺类型：0=淘宝 1=天猫
	shopType := ""
	if basic != nil {
		if ut, ok := basic["user_type"]; ok && ut != nil {
			switch fmt.Sprintf("%v", ut) {
			case "0":
				shopType = "淘宝"
			case "1":
				shopType = "天猫"
			default:
				shopType = fmt.Sprintf("%v", ut)
			}
		}
	}

	// 从 final_promotion_path_list.final_promotion_path_map_data 提取第一条优惠券信息
	couponType, couponFee, couponDesc, couponStart, couponEnd := extractFirstCoupon(price)

	return []interface{}{
		// A 列：adzone_id 留空（分组行已写）
		"",
		// 商品ID
		getString(item, "item_id"),
		// 标题
		getNestedString(basic, "title"),
		// 短标题
		getNestedString(basic, "short_title"),
		// 子标题
		getNestedString(basic, "sub_title"),
		// 品牌
		getNestedString(basic, "brand_name"),
		// 叶子分类
		getNestedString(basic, "category_name"),
		// 一级分类
		getNestedString(basic, "level_one_category_name"),
		// 店铺类型
		shopType,
		// 店铺名称
		getNestedString(basic, "shop_title"),
		// 月销量
		getNestedString(basic, "volume"),
		// 年销量
		getNestedString(basic, "annual_vol"),
		// 运费
		getNestedString(basic, "real_post_fee"),
		// 原价
		getNestedString(price, "reserve_price"),
		// 折扣价
		getNestedString(price, "zk_final_price"),
		// 最终成交价
		getNestedString(price, "final_promotion_price"),
		// 券类型
		couponType,
		// 优惠券面额
		couponFee,
		// 优惠券使用条件
		couponDesc,
		// 券开始时间
		couponStart,
		// 券结束时间
		couponEnd,
		// 凑单价
		getNestedString(price, "predict_rounding_up_price"),
		// 凑单描述
		getNestedString(price, "predict_rounding_up_price_desc"),
		// 佣金率
		getNestedString(publish, "income_rate"),
		// 主图
		addHttpsPrefix(getNestedString(basic, "pict_url")),
		// 推广链接
		addHttpsPrefix(getNestedString(publish, "click_url")),
		// 优惠券推广链接
		addHttpsPrefix(getNestedString(publish, "coupon_share_url")),
	}
}

// extractFirstCoupon 从 price_promotion_info 中提取 final_promotion_path_map_data 第一条促销信息
// 返回：促销类型、面额、使用条件、开始时间、结束时间
func extractFirstCoupon(price map[string]interface{}) (couponType, fee, desc, startTime, endTime string) {
	if price == nil {
		return
	}
	pathList, _ := price["final_promotion_path_list"].(map[string]interface{})
	if pathList == nil {
		return
	}
	mapData, _ := pathList["final_promotion_path_map_data"].([]interface{})
	if len(mapData) == 0 {
		return
	}
	first, _ := mapData[0].(map[string]interface{})
	if first == nil {
		return
	}
	couponType = fmt.Sprintf("%v", first["promotion_title"])
	fee = fmt.Sprintf("%v", first["promotion_fee"])
	desc = fmt.Sprintf("%v", first["promotion_desc"])
	startTime = msTimestampToDate(fmt.Sprintf("%v", first["promotion_start_time"]))
	endTime = msTimestampToDate(fmt.Sprintf("%v", first["promotion_end_time"]))
	return
}

// msTimestampToDate 将毫秒时间戳字符串转换为 "2006-01-02" 格式，失败时原样返回
func msTimestampToDate(ms string) string {
	if ms == "" || ms == "<nil>" {
		return ""
	}
	var ts int64
	if _, err := fmt.Sscanf(ms, "%d", &ts); err != nil {
		return ms
	}
	return time.Unix(ts/1000, 0).Format("2006-01-02")
}

// addHttpsPrefix 为缺少协议前缀的 URL 补全 https:
func addHttpsPrefix(u string) string {
	if u == "" {
		return ""
	}
	if strings.HasPrefix(u, "//") {
		return "https:" + u
	}
	return u
}

// getNestedString 从嵌套 map 中安全取字符串值
func getNestedString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
