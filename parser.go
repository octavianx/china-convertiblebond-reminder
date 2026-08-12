package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// BondParser 去掉东方财富返回的 JSONP 前缀(`_(...)`)后解析为通用结构。
// 之前这里丢弃了 Decode 的错误,导致数据源返回错误页(非 JSONP)时静默得到 nil;
// 现在把错误返回,让上层能区分「解析失败」与「今天没有可转债」。
func BondParser(data []byte) (any, error) {
	var bond any
	err := json.NewDecoder(&JsonpWrapper{
		Underlying: bytes.NewBuffer(data),
		Prefix:     "_",
	}).Decode(&bond)
	if err != nil {
		return nil, fmt.Errorf("解析 JSONP/JSON 失败: %w", err)
	}
	return bond, nil
}

// BondFilter 从解析结果中筛出今天/明天/后天可申购或预约的可转债。
// 返回 error 表示数据结构无法解析(与「今天没有可转债」区分开)。
func BondFilter(data any) (string, error) {
	type Bonds struct {
		Result struct {
			Data []struct {
				Name   string `json:"SECURITY_NAME_ABBR"`
				Code   string `json:"SECURITY_CODE"`
				Date   string `json:"VALUE_DATE"`
				Rating string `json:"RATING"`
			} `json:"data"`
		} `json:"result"`
	}
	var (
		message string = ""
		bonds   Bonds
	)

	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("重新编码数据失败: %w", err)
	}
	if err := json.Unmarshal(b, &bonds); err != nil {
		return "", fmt.Errorf("解析可转债数据结构失败: %w", err)
	}

	for _, v := range bonds.Result.Data {
		// 匹配今天
		if v.Date == time.Now().Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + "）\n"
		}
		// 匹配明天
		if v.Date == time.Now().Add(time.Hour*24).Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + " / 预约）\n"
		}
		// 匹配后天
		if v.Date == time.Now().Add(time.Hour*24*2).Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + " / 预约）\n"
		}
	}

	if len(message) == 0 {
		message = "今天没有可转债供申购或预约"
	}
	return message, nil
}
