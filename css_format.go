package main

import (
	"bytes"
	"log"
	"os"
)

func cssFormat(inPath string) []byte {
	// 读取整个文件
	data, err := os.ReadFile(inPath)
	if err != nil {
		log.Printf("读取文件(%s)失败: %v\n", inPath, err)
		return nil
	}

	// 统一换行符至 Unix 格式
	normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	normalized = bytes.ReplaceAll(normalized, []byte("\r"), []byte("\n"))

	// 拆分为行，并去除每行末尾多余空白
	lines := bytes.Split(normalized, []byte("\n"))
	for i := range lines {
		lines[i] = bytes.TrimRight(lines[i], " \t")
	}

	// 替换前导 Tab 为 2 个空格，保留原空格
	for i, line := range lines {
		// 处理前导空白
		var prefix []byte
		pos := 0
	scan:
		for j, b := range line {
			switch b {
			case '\t':
				// Tab 替换为两个空格
				prefix = append(prefix, ' ', ' ')
				pos = j + 1
			case ' ':
				// 原空格保留
				prefix = append(prefix, ' ')
				pos = j + 1
			default:
				// 遇到非空白字符停止
				break scan
			}
		}

		// 重建行：prefix + 剩余内容
		content := line[pos:]
		lines[i] = append(prefix, content...)
	}

	// 合并行，生成最终内容
	normalized = bytes.Join(lines, []byte("\n"))

	if bytes.Equal(data, normalized) {
		// 文件内容未改变，跳过写入
		return nil
	}
	return normalized
}
