package main

import (
	"encoding/json"
	"os"
)

// BarkTarget 是一个带名字的推送目标,便于区分设备、方便日志与 CRUD。
type BarkTarget struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Config struct {
	// Token / ChatId 是已废弃的 Telegram 字段,留待后续 telegram 清理分支移除。
	Token  string `json:"token"`
	ChatId int64  `json:"chat_id"`

	BarkUrl  string       `json:"barkurl"`   // 兼容旧的单目标写法
	BarkUrls []BarkTarget `json:"bark_urls"` // 多目标推送(推荐)
}

// Targets 汇总实际要推送的目标列表:
// 优先用 bark_urls;若为空则回退到旧的单个 barkurl。
func (config *Config) Targets() []BarkTarget {
	if len(config.BarkUrls) > 0 {
		return config.BarkUrls
	}
	if config.BarkUrl != "" {
		return []BarkTarget{{Name: "default", Url: config.BarkUrl}}
	}
	return nil
}

func (config *Config) ReadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}
