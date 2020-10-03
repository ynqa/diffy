package diffy

import (
	"strings"
)

type Option struct {
	NoHeader         bool
	TabSize          int
	SeparatorSymbol  string
	SeparatorWidth   int
	SpaceSizeAfterLn int
}

func countDigits(v int) int {
	var cnt int
	for v != 0 {
		v /= 10
		cnt++
	}
	return cnt
}

func formatTextLine(text string, tabSize int) string {
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "\t", strings.Repeat(" ", tabSize))
	return text
}

func splitText(text string, length, tabSize int) []string {
	text = formatTextLine(text, tabSize)
	if len(text) < length {
		return []string{text}
	}
	var res []string
	for i := 0; i < len(text); i += length {
		if i+length < len(text) {
			res = append(res, text[i:(i + length)])
		} else {
			res = append(res, text[i:])
		}
	}
	return res
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
