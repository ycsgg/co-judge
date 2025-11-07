package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"mips2hex/types"
)

func ReadFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines, s.Err()
}

// ParseLines 进行简单的两遍解析：收集 label 与 items (.word 也支持), .text 段启用
// .text 起始地址假定为 0
func ParseLines(lines []string) ([]types.Item, map[string]uint32, error) {
	items := []types.Item{}
	labels := map[string]uint32{}
	addr := uint32(0)
	inText := false
	lineNo := 0
	labelRe := regexp.MustCompile(`^([A-Za-z_\.][A-Za-z0-9_\.\$]*):`)

	for _, rawLine := range lines {
		lineNo++
		line := rawLine
		// 删除注释 (# 和 //)
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// directive
		if strings.HasPrefix(line, ".") {
			parts := fields(line)
			dir := strings.ToLower(parts[0])
			switch dir {
			case ".text":
				inText = true
				continue
			case ".data":
				inText = false
				continue
			case ".word":
				if !inText {
					continue
				}
				rest := strings.TrimSpace(line[len(".word"):])
				args := splitArgs(rest)
				for _, a := range args {
					it := types.Item{
						Kind:     types.Word,
						Raw:      a,
						Tokens:   []string{a},
						LineNo:   lineNo,
						OrigLine: rawLine,
					}
					items = append(items, it)
					addr += 4
				}
				continue
			default:
				// 忽略其他
				continue
			}
		}
		// labels
		for {
			m := labelRe.FindStringSubmatch(line)
			if m == nil {
				break
			}
			label := m[1]
			labels[label] = addr
			// 去掉该前缀 label: 部分，继续循环以处理多个 label
			line = strings.TrimSpace(line[len(m[0]):])
			if line == "" {
				// 如果 label 后面没有任何指令或内容，该行结束，继续到下一行
				break
			}
		}
		if line == "" {
			continue
		}
		if !inText {
			continue
		}
		toks := tokenize(line)
		if len(toks) == 0 {
			continue
		}
		it := types.Item{
			Kind:     types.Instr,
			Raw:      line,
			Tokens:   toks,
			LineNo:   lineNo,
			OrigLine: rawLine,
			Size:     4,
		}
		if strings.ToLower(toks[0]) == "li" {
			it.Size = 8
		}
		items = append(items, it)
		addr += it.Size
	}
	return items, labels, nil
}

func fields(s string) []string {
	return strings.Fields(s)
}

func splitArgs(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func tokenize(s string) []string {
	s = strings.ReplaceAll(s, ",", " ")
	parts := strings.Fields(s)
	return parts
}
