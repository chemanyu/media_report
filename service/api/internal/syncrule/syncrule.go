// Package syncrule 定义主生产 → 从节点同步的"部分同步"规则。
//
// 默认所有配置表都是整表覆盖（dump 全表 + 从节点 TRUNCATE 后重灌）。
// 部分同步解决两类问题：
//
//  1. 两边都在写的表：例如 fz_hourly_report，huawei 是外部 PHP 回传到主生产的
//     （从节点收不到），而 OPPO/小米/荣耀是从节点自己按凭据拉取的。整表覆盖
//     会把从节点自采的数据一起清掉，所以只能按 media 过滤同步。
//
//  2. 按日期无限累积的表：同上表，行数 = 天数 × 媒体数 × 账户数，且每 30 分钟
//     同步一次。全量同步会让传输量和事务大小随时间线性增长，而几个月前的历史
//     数据早已不再变动。用 IncrementalDays 划一个滚动日期窗口，只同步窗口内的
//     数据，窗口外的历史在从节点原地留存。
//
// 本包被 dump 接口（上游只导出匹配行）和 sync 调度器（下游只删除+重灌匹配行）
// 共同引用，保证两端规则永远一致 —— 规则写在代码里而不是 yaml 里是刻意的：
// deploy.sh 不推送 media-api.yaml，配在 yaml 里会导致两端各改一份、悄悄漂移，
// 而窗口不一致会让同步删掉不该删的数据。
package syncrule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Filter 描述一张表的部分同步规则。
type Filter struct {
	Column          string   // 值过滤列，例如 "media"；留空表示不按值过滤
	Values          []string // 允许同步的值，例如 ["huawei"]
	DateColumn      string   // 日期列，例如 "report_date"（int，格式 YYYYMMDD）
	IncrementalDays int      // 滚动窗口天数：只同步 DateColumn >= 今天-N 天的数据；0 表示不限（全部历史）
	DropColumns     []string // 写入从节点前需要剔除的列（通常是自增主键，避免与从节点自采数据主键冲突）
}

// Cond 是一个 SQL 条件片段，供两端拼到各自的查询上。
type Cond struct {
	SQL  string
	Args []any
}

// Conds 返回本规则对应的 SQL 条件。上游 dump 用它做 SELECT 的 WHERE，
// 下游同步用它做 DELETE 的 WHERE —— 同一个来源，保证删除范围和导出范围严格对齐。
// now 由调用方传入（通常是 time.Now()），便于测试。
func (f Filter) Conds(now time.Time) []Cond {
	var conds []Cond
	if f.Column != "" && len(f.Values) > 0 {
		conds = append(conds, Cond{
			SQL:  fmt.Sprintf("`%s` IN ?", f.Column),
			Args: []any{f.Values},
		})
	}
	if cutoff, ok := f.Cutoff(now); ok {
		conds = append(conds, Cond{
			SQL:  fmt.Sprintf("`%s` >= ?", f.DateColumn),
			Args: []any{cutoff},
		})
	}
	return conds
}

// Cutoff 返回滚动窗口的起始日期（YYYYMMDD 整数）；ok=false 表示该表不限日期。
// 注意窗口按调用方所在时区的"今天"计算，两端服务器时区需一致（均为 Asia/Shanghai）。
func (f Filter) Cutoff(now time.Time) (int, bool) {
	if f.DateColumn == "" || f.IncrementalDays <= 0 {
		return 0, false
	}
	d := now.AddDate(0, 0, -f.IncrementalDays)
	return d.Year()*10000 + int(d.Month())*100 + d.Day(), true
}

// Desc 返回规则的可读描述，用于日志。
func (f Filter) Desc(now time.Time) string {
	var parts []string
	if f.Column != "" && len(f.Values) > 0 {
		parts = append(parts, fmt.Sprintf("%s IN %v", f.Column, f.Values))
	}
	if cutoff, ok := f.Cutoff(now); ok {
		parts = append(parts, fmt.Sprintf("%s >= %d (最近 %d 天)", f.DateColumn, cutoff, f.IncrementalDays))
	}
	if len(parts) == 0 {
		return "全表"
	}
	return strings.Join(parts, " AND ")
}

// Match 判断一行数据是否属于本规则的同步范围（值命中 + 在日期窗口内）。
// 上游理论上已经过滤过，这里是下游防线：两端规则漂移时不至于把范围外的数据写进来。
func (f Filter) Match(row map[string]any, now time.Time) bool {
	if f.Column != "" && len(f.Values) > 0 {
		got, ok := toString(row[f.Column])
		if !ok {
			return false
		}
		matched := false
		for _, v := range f.Values {
			if got == v {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if cutoff, ok := f.Cutoff(now); ok {
		date, valid := toInt(row[f.DateColumn])
		if !valid || date < cutoff {
			return false
		}
	}
	return true
}

var filters = map[string]Filter{
	// huawei 数据由外部 PHP 回传到主生产，从节点收不到，需要从主生产同步；
	// 同表其他媒体（OPPO/小米/荣耀）从节点自己拉取，不参与同步。
	// 只滚动同步最近 30 天：回传报表过了几天就不再修订，更早的历史无需反复传输。
	"fz_hourly_report": {
		Column:          "media",
		Values:          []string{"huawei"},
		DateColumn:      "report_date",
		IncrementalDays: 7,
		DropColumns:     []string{"id"},
	},
}

// For 返回指定表的部分同步规则；ok=false 表示该表整表覆盖。
func For(table string) (Filter, bool) {
	f, ok := filters[table]
	return f, ok
}

func toString(v any) (string, bool) {
	switch s := v.(type) {
	case string:
		return s, true
	case []byte:
		return string(s), true
	default:
		return "", false
	}
}

// toInt 兼容各种数值来源：本地 GORM 查询返回 int64，
// 经 JSON 往返后变成 float64，某些驱动/列类型下还可能是 []byte / string。
func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case uint64:
		return int(n), true
	case float64:
		return int(n), true
	case float32:
		return int(n), true
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(n))
		return i, err == nil
	case []byte:
		i, err := strconv.Atoi(strings.TrimSpace(string(n)))
		return i, err == nil
	default:
		return 0, false
	}
}
