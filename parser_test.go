package main

import (
	"strings"
	"testing"
	"time"
)

// 数据源返回错误页(非 JSONP)时,必须报错,而不是静默得到 nil。
func TestBondParser_NonJSONP_ReturnsError(t *testing.T) {
	_, err := BondParser([]byte("<html><body>403 Forbidden</body></html>"))
	if err == nil {
		t.Fatal("非 JSONP 输入应返回错误,却返回 nil error(会被伪装成『今天没有』)")
	}
}

// JSONP 前缀正确但 JSON 损坏时,同样必须报错。
func TestBondParser_BrokenJSON_ReturnsError(t *testing.T) {
	_, err := BondParser([]byte("_({ broken json "))
	if err == nil {
		t.Fatal("损坏 JSON 应返回错误")
	}
}

// 合法但今日无债:返回「今天没有」提示,且不报错。
func TestBondFilter_EmptyIsNotError(t *testing.T) {
	parsed, err := BondParser([]byte(`_({"result":{"data":[]}})`))
	if err != nil {
		t.Fatalf("合法空数据不应报错: %v", err)
	}
	msg, err := BondFilter(parsed)
	if err != nil {
		t.Fatalf("BondFilter 不应报错: %v", err)
	}
	if !strings.Contains(msg, "今天没有") {
		t.Fatalf("空数据应提示『今天没有』,实际: %q", msg)
	}
}

// 命中今日可申购:消息里应出现债券名与代码。
func TestBondFilter_TodayMatch(t *testing.T) {
	today := time.Now().Format("2006-01-02") + " 00:00:00"
	payload := `_({"result":{"data":[{"SECURITY_NAME_ABBR":"测试转债","SECURITY_CODE":"123456","VALUE_DATE":"` + today + `","RATING":"AA"}]}})`

	parsed, err := BondParser([]byte(payload))
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	msg, err := BondFilter(parsed)
	if err != nil {
		t.Fatalf("筛选失败: %v", err)
	}
	if !strings.Contains(msg, "测试转债") || !strings.Contains(msg, "123456") {
		t.Fatalf("应包含债券名与代码,实际: %q", msg)
	}
	if strings.Contains(msg, "今天没有") {
		t.Fatalf("命中时不应出现『今天没有』,实际: %q", msg)
	}
}
