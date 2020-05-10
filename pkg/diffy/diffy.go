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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
