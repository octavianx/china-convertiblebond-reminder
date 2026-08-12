package main

import (
	"bytes"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageSender2 向单个 Bark 目标推送一条消息。
// 目标(名字 + URL)由 config.json 提供,不再硬编码;
// 多目标推送时由上层对每个 target 各调用一次。
func MessageSender2(target BarkTarget, message string) {
	if target.Url == "" {
		log.Printf("目标 [%s] 的 bark url 为空,跳过推送", target.Name)
		return
	}

	log.Printf("sending to [%s]: %s", target.Name, target.Url)

	body := "body=" + message
	req, err := http.NewRequest("POST", target.Url, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("[%s] Error creating request: %v", target.Name, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%s] Error performing request: %v", target.Name, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[%s] Response Status: %s", target.Name, resp.Status)
}

func MessageSender(bot *tgbotapi.BotAPI, id int64, message string) {
	// Telegram 通道已废弃,保留空实现以兼容旧签名。
	//	bot.Send(
	//		tgbotapi.NewMessage(
	//			id, message,
	//		),
	//	)

	log.Println("Dummy Message sent successfully")
}
