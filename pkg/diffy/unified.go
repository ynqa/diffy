package diffy

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/gookit/color.v1"
)

func WriteUnifiedDiff(
	w io.Writer,
	diff difflib.UnifiedDiff,
	opt Option,
) error {
	lnSpaceSize := countDigits(max(len(diff.A), len(diff.B)))

	buf := bufio.NewWriter(w)
	defer buf.Flush()

	groupedOpcodes := difflib.NewMatcher(diff.A, diff.B).GetGroupedOpCodes(diff.Context)
	for i, opcodes := range groupedOpcodes {
		if i == 0 && !opt.NoHeader {
			buf.WriteString(coloredUnifiedHeader(diff.FromFile, diff.ToFile))
		}
		for _, c := range opcodes {
			i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
			if c.Tag == 'e' {
				for ln, line := range diff.A[i1:i2] {
					buf.WriteString(
						fmt.Sprintf(
							"%s %s%s%s\n",
							color.New(color.Gray).Sprintf("%*d", lnSpaceSize, i1+ln+1),
							color.New(color.Gray).Sprintf("%*d", lnSpaceSize, j1+ln+1),
							strings.Repeat(" ", opt.SpaceSizeAfterLn),
							fmt.Sprintf("%s", formatTextLine(line, opt.TabSize)),
						),
					)
				}
			}
			if c.Tag == 'r' || c.Tag == 'd' {
				for ln, line := range diff.A[i1:i2] {
					buf.WriteString(
						fmt.Sprintf(
							" %s%s%s\n",
							color.New(color.Gray).Sprintf("%*d", lnSpaceSize, i1+ln+1),
							strings.Repeat(" ", lnSpaceSize+opt.SpaceSizeAfterLn),
							color.New(color.Red, color.Bold).Sprintf("%s", formatTextLine(line, opt.TabSize)),
						),
					)
				}
			}
			if c.Tag == 'r' || c.Tag == 'i' {
				for ln, line := range diff.B[j1:j2] {
					buf.WriteString(
						fmt.Sprintf(
							" %s%s%s\n",
							color.New(color.Gray).Sprintf("%*d", lnSpaceSize*2, j1+ln+1),
							strings.Repeat(" ", opt.SpaceSizeAfterLn),
							color.New(color.Green, color.Bold).Sprintf("%s", formatTextLine(line, opt.TabSize)),
						),
					)
				}
			}
		}
		if i != len(groupedOpcodes)-1 {
			buf.WriteString(color.New(color.Blue).Sprintf("%s\n", strings.Repeat(opt.SeparatorSymbol, opt.SeparatorWidth)))
		}
	}
	return nil
}

func coloredUnifiedHeader(org, new string) string {
	return fmt.Sprintf(
		"%s %s\n%s %s\n",
		color.New(color.Red).Sprint("-"),
		color.New(color.Bold).Sprint(org),
		color.New(color.Green).Sprint("+"),
		color.New(color.Bold).Sprint(new),
	)
}
