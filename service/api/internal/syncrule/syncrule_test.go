package syncrule

import (
	"testing"
	"time"
)

// 固定基准时间，避免测试依赖真实时钟
var base = time.Date(2026, 7, 29, 10, 0, 0, 0, time.Local)

func TestCutoff(t *testing.T) {
	f, ok := For("fz_hourly_report")
	if !ok {
		t.Fatal("fz_hourly_report 规则丢失")
	}

	cutoff, ok := f.Cutoff(base)
	if !ok {
		t.Fatal("期望有日期窗口")
	}
	// 2026-07-29 往前 30 天 = 2026-06-29
	if cutoff != 20260629 {
		t.Errorf("cutoff = %d, 期望 20260629", cutoff)
	}

	// 跨月/跨年边界：AddDate 应正确回退
	cutoff, _ = f.Cutoff(time.Date(2026, 1, 10, 0, 0, 0, 0, time.Local))
	if cutoff != 20251211 {
		t.Errorf("跨年 cutoff = %d, 期望 20251211", cutoff)
	}

	// 未配日期列的规则不应产生窗口
	var noDate Filter
	if _, ok := noDate.Cutoff(base); ok {
		t.Error("未配 DateColumn 时不应有窗口")
	}
}

func TestMatch(t *testing.T) {
	f, _ := For("fz_hourly_report")

	cases := []struct {
		name string
		row  map[string]any
		want bool
	}{
		{"窗口内的 huawei", map[string]any{"media": "huawei", "report_date": 20260728}, true},
		{"恰好等于 cutoff", map[string]any{"media": "huawei", "report_date": 20260629}, true},
		{"cutoff 前一天", map[string]any{"media": "huawei", "report_date": 20260628}, false},
		{"窗口内但媒体不符", map[string]any{"media": "oppo", "report_date": 20260728}, false},
		{"JSON 往返后 float64", map[string]any{"media": "huawei", "report_date": float64(20260728)}, true},
		{"驱动返回 []byte", map[string]any{"media": []byte("huawei"), "report_date": []byte("20260728")}, true},
		{"日期列缺失", map[string]any{"media": "huawei"}, false},
		{"媒体列缺失", map[string]any{"report_date": 20260728}, false},
	}
	for _, c := range cases {
		if got := f.Match(c.row, base); got != c.want {
			t.Errorf("%s: Match = %v, 期望 %v", c.name, got, c.want)
		}
	}
}

func TestConds(t *testing.T) {
	f, _ := For("fz_hourly_report")
	conds := f.Conds(base)
	if len(conds) != 2 {
		t.Fatalf("期望 2 个条件（媒体 + 日期），得到 %d", len(conds))
	}
	if conds[0].SQL != "`media` IN ?" {
		t.Errorf("条件[0] = %q", conds[0].SQL)
	}
	if conds[1].SQL != "`report_date` >= ?" {
		t.Errorf("条件[1] = %q", conds[1].SQL)
	}
	if got := conds[1].Args[0]; got != 20260629 {
		t.Errorf("日期参数 = %v, 期望 20260629", got)
	}
}

// 无规则的表必须走整表覆盖，不能被误判成部分同步（否则 DELETE 范围会出错）
func TestForUnconfiguredTable(t *testing.T) {
	for _, tbl := range []string{"media_token", "fz_config", "qczj_report_data"} {
		if _, ok := For(tbl); ok {
			t.Errorf("%s 不应有部分同步规则", tbl)
		}
	}
}
