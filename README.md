# 1. simple-formater

### 1.1 使用 Tab 符号现代化 CSS/HTML/JS 代码格式
1. 统一前导空白字符（空格和 Tab）为 Tab：  
将每行代码开头的所有前导空白字符（空格或 Tab）转换为 Tab 字符（\t）。
2. 移除每行末尾的空白字符（空格和 Tab） ：
删除每行末尾的所有空白字符（空格或 Tab）。

#### 1.1.1 示例：
处理前：
```txt
    // 4 个空格  
→     // 1 个 Tab + 2 个空格  
Hello World     ← 末尾空格  
```
处理后：
```txt
→   // 统一为 Tab  
→   →   // 仅保留 Tab  
Hello World ← 末尾空格已删除  
```

### 1.2 Modernize code presentation with tab for CSS/HTML/JS formatting
1. Normalize leading whitespace (spaces and tabs) to tabs:  
Convert all leading whitespace characters (spaces or tabs) at the start of each line to tab characters (\t).
2. Remove trailing whitespace (spaces and tabs) from each line:  
Trim any trailing whitespace characters (spaces or tabs) at the end of each line.

#### 1.2.1 Example:
Before:
```txt
    // 4 spaces
→     // 1 tab + 2 spaces
Hello World     ← Trailing spaces
```
After:
```txt
→   // Normalized to tabs
→   →   // Tabs only
Hello World ← Trailing spaces removed
```

# 2. usage

```sh
用法：./simple-formater -dir <directory>

选项：
  -dir string
        要遍历的根目录路径（必填）
  -h    显示帮助信息（等同于 -help）
  -help
        显示帮助信息

<bug report>    https://github.com/wiselike/simple-formater/issues
```
