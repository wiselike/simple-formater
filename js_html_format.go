package main

import (
	"bytes"
	"log"
	"os"
)

func jsHtmlFormat(inPath string) []byte {
	// 读取原始文件（按字节读取，保留中文字符）
	original, err := os.ReadFile(inPath)
	if err != nil {
		log.Printf("读取文件(%s)失败: %v\n", inPath, err)
		return nil
	}

	// 将 Windows 与旧 Mac 换行符统一替换为 Unix 格式 (\n)
	normalized := bytes.ReplaceAll(original, []byte("\r\n"), []byte("\n"))
	normalized = bytes.ReplaceAll(normalized, []byte("\r"), []byte("\n"))

	// 拆分为行，并去除每行末尾多余空白
	lines := bytes.Split(normalized, []byte("\n"))
	for i := range lines {
		lines[i] = bytes.TrimRight(lines[i], " \t")
	}

	// 初始化前一行修改后缩进量和有效缩进
	prevModifiedTabs := 0
	prevEffective := 0

	for i := range lines {
		// 跳过空行或仅含空白的行
		if len(lines[i]) == 0 {
			continue
		}

		// 统计当前行开头的空格和制表符，并计算有效缩进（Tab 对齐到下一个 4 空格边界）
		countSpaces, countTabs, effective := 0, 0, 0
	scan:
		for _, b := range lines[i] {
			switch b {
			case ' ':
				countSpaces++ // 空格计数
				effective++   // 有效缩进加一
			case '\t':
				countTabs++ // Tab 计数
				// 有效缩进跳到下一个 4 的倍数边界
				effective = ((effective / 4) + 1) * 4
			default:
				break scan
			}
		}
		currEffective := effective
		newTabs := prevModifiedTabs

		// 若行前导全为 Tab，则保留原始 Tab 数
		if countSpaces == 0 && countTabs > 0 {
			newTabs = countTabs
		} else {
			// 比较当前有效缩进与上一行有效缩进
			if currEffective > prevEffective {
				newTabs = prevModifiedTabs + 1
			} else if currEffective < prevEffective {
				newTabs = prevModifiedTabs - 1
			}
			if newTabs < 0 {
				newTabs = 0
			}
		}

		// 清除原有前导空白，仅保留内容，并应用 newTabs
		content := bytes.TrimLeft(lines[i], " \t")
		lines[i] = append(bytes.Repeat([]byte("\t"), newTabs), content...)

		// 更新上一行信息
		prevModifiedTabs = newTabs
		prevEffective = currEffective
	}

	// 合并内容
	normalized = bytes.Join(lines, []byte("\n"))

	// 比较内容变化，决定输出信息
	if bytes.Equal(original, normalized) {
		// 文件内容未改变，跳过写入
		return nil
	}
	return normalized
}
