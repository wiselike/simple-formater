package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// 自定义帮助信息
func printHelp() {
	fmt.Fprintf(flag.CommandLine.Output(), "用法：%s -dir <directory>\n\n", os.Args[0])
	fmt.Println("选项：")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("<bug report>\thttps://github.com/wiselike/simple-formater/issues")
}

// readConfig 会尝试在 dir 目录下读取 .formatrc 文件，
// 如果文件不存在则返回 disable=false、nil 错误；
// 如果能读到并解析成功，返回文件中 disable 的值。
func readConfig(dir string) (disable bool) {
	configPath := filepath.Join(dir, ".formatrc")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}
	var cfg struct {
		Disable bool `json:"disable"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return false
	}
	return cfg.Disable
}

func main() {
	// 定义命令行参数
	dir := flag.String("dir", "", "要遍历的根目录路径（必填）")
	help := flag.Bool("help", false, "显示帮助信息")
	flag.BoolVar(help, "h", false, "显示帮助信息（等同于 -help）")

	flag.Usage = printHelp
	flag.Parse()

	// 如果请求帮助，则打印并退出
	if *help {
		flag.Usage()
		return
	}

	// 校验必填参数
	if *dir == "" {
		fmt.Println("错误：必须指定 -dir 参数")
		flag.Usage()
		os.Exit(1)
	}

	// 日志时间格式
	log.SetFlags(log.LstdFlags)

	// 并发控制：CPU 核心数 * 5
	maxGoroutines := runtime.NumCPU() * 5
	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	// 遍历目录
	err := filepath.WalkDir(*dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("遍历 %s 时出错: %v", path, err)
			return err
		}

		// 1. 目录优先：检查 .formatrc，并根据 disable 决定是否跳过整个目录
		if d.IsDir() {
			if disable := readConfig(path); disable {
				return filepath.SkipDir
			}
			return nil
		}

		// 2. 如果是配置文件本身，跳过处理
		if d.Name() == ".formatrc" {
			return nil
		}

		// 3. 普通文件：并发执行 work
		wg.Add(1)
		sem <- struct{}{} // 获取令牌
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放令牌

			work(p)
		}(path)

		return nil
	})
	if err != nil {
		log.Fatalf("遍历目录时出错: %v", err)
	}

	wg.Wait()
}

// work 是要在每个文件路径上执行的函数
func work(path string) {
	var normalized []byte

	switch filepath.Ext(path) {
	case ".css":
		normalized = cssFormat(path)
	case ".js", ".html":
		normalized = jsHtmlFormat(path)
	}

	if normalized != nil {
		err := os.WriteFile(path, normalized, 0644)
		if err != nil {
			log.Printf("写入文件(%s)失败：%v\n", path, err)
			return
		}
		fmt.Println(path)
	}
}
