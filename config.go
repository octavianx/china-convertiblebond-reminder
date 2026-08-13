package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// BarkTarget 是一个带名字的推送目标,便于区分设备、方便日志与 CRUD。
type BarkTarget struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Device 是一台设备:引用某个服务器(servers 里的名字,或直接写 URL)+ 该设备的 key。
type Device struct {
	Name   string `json:"name"`
	Server string `json:"server"` // servers 注册表里的名字,或字面 http(s) URL
	Key    string `json:"key"`
}

type Config struct {
	// Token / ChatId 是已废弃的 Telegram 字段,留待后续 telegram 清理分支移除。
	Token  string `json:"token"`
	ChatId int64  `json:"chat_id"`

	// 多服务器 + 多设备(推荐):servers 定义一次,devices 引用。
	Servers map[string]string `json:"servers"`
	Devices []Device          `json:"devices"`

	BarkUrl  string       `json:"barkurl"`   // 兼容:旧的单目标写法
	BarkUrls []BarkTarget `json:"bark_urls"` // 兼容:每设备一条完整 URL
}

// Targets 汇总实际要推送的目标列表。
// 优先级:devices(servers 解析)> bark_urls > 单个 barkurl。
func (config *Config) Targets() []BarkTarget {
	if len(config.Devices) > 0 {
		var out []BarkTarget
		for _, d := range config.Devices {
			base := config.resolveServer(d.Server)
			if base == "" {
				log.Printf("设备 [%s] 的 server %q 既不在 servers 注册表、也不是 URL,跳过", d.Name, d.Server)
				continue
			}
			if d.Key == "" {
				log.Printf("设备 [%s] 缺 key,跳过", d.Name)
				continue
			}
			out = append(out, BarkTarget{Name: d.Name, Url: joinURL(base, d.Key)})
		}
		return out
	}
	if len(config.BarkUrls) > 0 {
		return config.BarkUrls
	}
	if config.BarkUrl != "" {
		return []BarkTarget{{Name: "default", Url: config.BarkUrl}}
	}
	return nil
}

// resolveServer 把 device.server 解析成服务器 base URL:
// 先查 servers 注册表;查不到则允许它本身就是一个 http(s) URL(字面写法)。
func (config *Config) resolveServer(server string) string {
	if v, ok := config.Servers[server]; ok {
		return v
	}
	if strings.HasPrefix(server, "http://") || strings.HasPrefix(server, "https://") {
		return server
	}
	return ""
}

// joinURL 拼 base + key,容忍 base 末尾的斜杠。
func joinURL(base, key string) string {
	return strings.TrimRight(base, "/") + "/" + key
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
