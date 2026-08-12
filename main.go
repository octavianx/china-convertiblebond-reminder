package main

import (
	"log"
)

// buildBondMessage 抓取并解析可转债数据,返回要推送的文本。
// 任何一步失败都返回带 ⚠️ 的明确告警,而不是伪装成「今天没有可转债」——
// 后者会与真实的「无新债」无法区分,导致数据源故障时被悄悄忽略。
func buildBondMessage() string {
	raw, err := BondData()
	if err != nil {
		log.Println("抓取失败:", err)
		return "⚠️ 可转债提醒:数据抓取失败,请检查网络或数据源\n" + err.Error()
	}

	parsed, err := BondParser(raw)
	if err != nil {
		log.Println("解析失败:", err)
		return "⚠️ 可转债提醒:数据解析失败,数据源格式可能已变更\n" + err.Error()
	}

	message, err := BondFilter(parsed)
	if err != nil {
		log.Println("筛选失败:", err)
		return "⚠️ 可转债提醒:数据结构解析失败,请检查数据源\n" + err.Error()
	}

	return message
}

func main() {
	var (
		args Args
		conf Config
	)

	args.ReadFlags()
	if err := conf.ReadConfig(args.Path); err != nil {
		log.Fatalln("读取 config 失败:", err)
	}

	// 调度由外部 cron 决定(工作日 09:05 / 09:35,见 docs/SPEC.md):
	// 本程序每次运行只推一次即退出,不自行判断时间或工作日。
	log.Println("Sending data")
	message := buildBondMessage()
	log.Println("Got feedback")

	targets := conf.Targets()
	if len(targets) == 0 {
		log.Println("未配置任何 bark 推送目标,跳过")
	}
	for _, t := range targets {
		MessageSender2(t, message)
	}
}
