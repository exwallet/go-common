package configuration

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadJSONFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	results := ""
	// 逐行读取
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		line = strings.TrimSpace(line)
		// 注释行
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		// 行中有//注释
		if strings.Index(line, "//") >= 0 {
			rs := []rune(line)
			index := strings.Index(line, "//")
			line = string(rs[:index])
		}
		// 行中有#注释
		if strings.Index(line, "#") >= 0 {
			rs := []rune(line)
			index := strings.Index(line, "#")
			line = string(rs[:index])
		}
		results += line
	}
	return results
}
