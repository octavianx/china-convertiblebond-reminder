package main

import (
	"flag"
	"os"
	"path/filepath"
)

type Args struct {
	Path    string
	Version bool
}

func (args *Args) ReadFlags() {
	flag.StringVar(&args.Path, "config", "", "config 路径(默认:二进制同目录下的 config.json)")
	flag.BoolVar(&args.Version, "version", false, "Show version information")
	flag.Parse()

	// 未显式指定时,默认在二进制自身所在目录找 config.json,
	// 而不是依赖当前工作目录(CWD)——这样 cron 从任意目录拉起都能找到。
	if args.Path == "" {
		args.Path = defaultConfigPath()
	}
}

// defaultConfigPath 返回二进制同目录下的 config.json;
// 若无法定位可执行文件,回退到 CWD 下的 ./config.json。
func defaultConfigPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(filepath.Dir(exe), "config.json")
}
