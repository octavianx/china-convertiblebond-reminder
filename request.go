package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const bondAPIURL = "https://datacenter-web.eastmoney.com/api/data/v1/get?callback=_&sortColumns=PUBLIC_START_DATE&sortTypes=-1&pageNumber=1&quoteType=0&reportName=RPT_BOND_CB_LIST&columns=ALL&quoteColumns=f2~01~CONVERT_STOCK_CODE~CONVERT_STOCK_PRICE,f235~10~SECURITY_CODE~TRANSFER_PRICE,f236~10~SECURITY_CODE~TRANSFER_VALUE,f2~10~SECURITY_CODE~CURRENT_BOND_PRICE,f237~10~SECURITY_CODE~TRANSFER_PREMIUM_RATIO,f239~10~SECURITY_CODE~RESALE_TRIG_PRICE,f240~10~SECURITY_CODE~REDEEM_TRIG_PRICE,f23~01~CONVERT_STOCK_CODE~PBV_RATIO"

// BondData 抓取可转债列表原始字节。
// 最多重试 maxRetries 次;全部失败(网络错误 / 非 200 响应 / 读取失败)时返回 error,
// 交由上层推送明确的失败告警,而不是无限重试或静默返回空数据。
func BondData() ([]byte, error) {
	const maxRetries = 5

	client := http.Client{Timeout: 10 * time.Second}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("抓取可转债数据(第 %d/%d 次)", attempt, maxRetries)

		resp, err := client.Get(bondAPIURL)
		if err != nil {
			lastErr = fmt.Errorf("网络请求失败: %w", err)
			log.Println(lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("数据源返回非 200 状态: %s", resp.Status)
			log.Println(lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		b, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("读取响应体失败: %w", err)
			log.Println(lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		return b, nil
	}

	return nil, fmt.Errorf("连续 %d 次抓取失败,放弃: %w", maxRetries, lastErr)
}
