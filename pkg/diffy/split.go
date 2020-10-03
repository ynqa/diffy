package diffy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/gookit/color.v1"
)

func WriteSplitDiff(
	w io.Writer,
	diff difflib.UnifiedDiff,
	opt Option,
) error {
	lnSpaceSize := countDigits(max(len(diff.A), len(diff.B)))
	width, _, err := terminalShape()
	if err != nil {
		return err
	}
	mid := width / 2

	buf := bufio.NewWriter(w)
	defer buf.Flush()

	groupedOpcodes := difflib.NewMatcher(diff.A, diff.B).GetGroupedOpCodes(diff.Context)
	for i, opcodes := range groupedOpcodes {
		if i == 0 && !opt.NoHeader {
			buf.WriteString(splittedHeader(diff.FromFile, diff.ToFile, color.New(color.Bold), mid))
		}
		for _, c := range opcodes {
			i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
			if c.Tag == 'e' {
				for ln, line := range diff.A[i1:i2] {
					texts := splitText(line, mid-2-lnSpaceSize*2, opt.TabSize)
					buf.WriteString(
						splittedLine(
							fmt.Sprintf("%*d", lnSpaceSize, i1+ln+1),
							texts[0],
							fmt.Sprintf("%*d", lnSpaceSize, j1+ln+1),
							texts[0],
							color.New(color.Gray),
							color.New(),
							color.New(),
							mid, opt.SpaceSizeAfterLn, 2,
						),
					)
					for i := 1; i < len(texts); i++ {
						buf.WriteString(
							splittedLine(
								strings.Repeat(" ", lnSpaceSize),
								texts[i],
								strings.Repeat(" ", lnSpaceSize),
								texts[i],
								color.New(color.Gray),
								color.New(),
								color.New(),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
					}
				}
			}
			if c.Tag == 'd' {
				for ln, line := range diff.A[i1:i2] {
					texts := splitText(line, mid-2-lnSpaceSize*2, opt.TabSize)
					buf.WriteString(
						splittedLine(
							fmt.Sprintf("%*d", lnSpaceSize, i1+ln+1),
							texts[0],
							"",
							"",
							color.New(color.Gray),
							color.New(color.Red, color.Bold),
							color.New(),
							mid, opt.SpaceSizeAfterLn, 2,
						),
					)
					for i := 1; i < len(texts); i++ {
						buf.WriteString(
							splittedLine(
								strings.Repeat(" ", lnSpaceSize),
								texts[i],
								"",
								"",
								color.New(color.Gray),
								color.New(color.Red, color.Bold),
								color.New(),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
					}
				}
			}
			if c.Tag == 'i' {
				for ln, line := range diff.B[j1:j2] {
					texts := splitText(line, mid-2-lnSpaceSize*2, opt.TabSize)
					buf.WriteString(
						splittedLine(
							"",
							"",
							fmt.Sprintf("%*d", lnSpaceSize, j1+ln+1),
							texts[0],
							color.New(color.Gray),
							color.New(),
							color.New(color.Green, color.Bold),
							mid, opt.SpaceSizeAfterLn, 2,
						),
					)
					for i := 1; i < len(texts); i++ {
						buf.WriteString(
							splittedLine(
								"",
								"",
								strings.Repeat(" ", lnSpaceSize),
								texts[i],
								color.New(color.Gray),
								color.New(),
								color.New(color.Green, color.Bold),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
					}
				}
			}
			if c.Tag == 'r' {
				cursor := 0
				// compare the number of raws
				minLen := i2 - i1
				minIsOrg := true
				if j2-j1 < minLen {
					minLen = j2 - j1
					minIsOrg = false
				}
				for ; cursor < minLen; cursor++ {
					orgTexts := splitText(diff.A[i1+cursor], mid-2-lnSpaceSize*2, opt.TabSize)
					newTexts := splitText(diff.B[j1+cursor], mid-2-lnSpaceSize*2, opt.TabSize)
					buf.WriteString(
						splittedLine(
							fmt.Sprintf("%*d", lnSpaceSize, i1+cursor+1),
							orgTexts[0],
							fmt.Sprintf("%*d", lnSpaceSize, j1+cursor+1),
							newTexts[0],
							color.New(color.Gray),
							color.New(color.Red, color.Bold),
							color.New(color.Green, color.Bold),
							mid, opt.SpaceSizeAfterLn, 2,
						),
					)
					textCursor := 1
					minTextLen := len(orgTexts)
					minTextIsOrg := true
					if len(newTexts) < minTextLen {
						minTextLen = len(newTexts)
						minTextIsOrg = false
					}
					for ; textCursor < minTextLen; textCursor++ {
						buf.WriteString(
							splittedLine(
								strings.Repeat(" ", lnSpaceSize),
								orgTexts[textCursor],
								strings.Repeat(" ", lnSpaceSize),
								newTexts[textCursor],
								color.New(color.Gray),
								color.New(color.Red, color.Bold),
								color.New(color.Green, color.Bold),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
					}
					if minTextIsOrg {
						for i := textCursor; i < len(newTexts); i++ {
							buf.WriteString(
								splittedLine(
									"",
									"",
									strings.Repeat(" ", lnSpaceSize),
									newTexts[i],
									color.New(color.Gray),
									color.New(),
									color.New(color.Green, color.Bold),
									mid, opt.SpaceSizeAfterLn, 2,
								),
							)
						}
					} else {
						for i := textCursor; i < len(orgTexts); i++ {
							buf.WriteString(
								splittedLine(
									strings.Repeat(" ", lnSpaceSize),
									orgTexts[i],
									"",
									"",
									color.New(color.Gray),
									color.New(color.Red, color.Bold),
									color.New(),
									mid, opt.SpaceSizeAfterLn, 2,
								),
							)
						}
					}
				}
				if minIsOrg {
					for ; cursor < j2-j1; cursor++ {
						texts := splitText(diff.B[j1+cursor], mid-2-lnSpaceSize*2, opt.TabSize)
						buf.WriteString(
							splittedLine(
								"",
								"",
								fmt.Sprintf("%*d", lnSpaceSize, j1+cursor+1),
								texts[0],
								color.New(color.Gray),
								color.New(),
								color.New(color.Green, color.Bold),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
						for i := 1; i < len(texts); i++ {
							buf.WriteString(
								splittedLine(
									"",
									"",
									strings.Repeat(" ", lnSpaceSize),
									texts[i],
									color.New(color.Gray),
									color.New(),
									color.New(color.Green, color.Bold),
									mid, opt.SpaceSizeAfterLn, 2,
								),
							)
						}
					}
				} else {
					for ; cursor < i2-i1; cursor++ {
						texts := splitText(diff.A[i1+cursor], mid-2-lnSpaceSize*2, opt.TabSize)
						buf.WriteString(
							splittedLine(
								fmt.Sprintf("%*d", lnSpaceSize, i1+cursor+1),
								texts[0],
								"",
								"",
								color.New(color.Gray),
								color.New(color.Red, color.Bold),
								color.New(),
								mid, opt.SpaceSizeAfterLn, 2,
							),
						)
						for i := 1; i < len(texts); i++ {
							buf.WriteString(
								splittedLine(
									strings.Repeat(" ", lnSpaceSize),
									texts[i],
									"",
									"",
									color.New(color.Gray),
									color.New(color.Red, color.Bold),
									color.New(),
									mid, opt.SpaceSizeAfterLn, 2,
								),
							)
						}
					}
				}
			}
		}
		if i != len(groupedOpcodes)-1 {
			buf.WriteString(
				splittedFootLine(
					opt.SeparatorSymbol,
					opt.SeparatorWidth,
					color.New(color.Blue),
					mid,
				),
			)
		}
	}
	return nil
}

func get(i int, strs []string) string {
	if i < len(strs) {
		return strs[i]
	}
	return ""
}

func rpad(line string, limit int) string {
	if len(line) < limit {
		line += strings.Repeat(" ", limit-len(line))
	}
	return line
}

func chopped(line string, limit int) []string {
	var res []string
	if limit >= len(line) {
		res = []string{line}
	} else {
		for limit < len(line) {
			res = append(res, line[:limit])
			line = line[limit:]
		}
		res = append(res, line)
	}
	return res
}

func splittedFootLine(symbol string, width int, style color.Style, boundary int) string {
	d := strings.Repeat(symbol, width)
	return style.Sprintf("%s%s\n", rpad(d, boundary), d)
}

func splittedHeader(rawLeftFile, rawRightFile string, style color.Style, boundary int) string {
	return style.Sprintf("%s%s\n", rpad(rawLeftFile, boundary), rawRightFile)
}

func splittedLine(
	rawLeftLn, rawLeft, rawRightLn, rawRight string,
	lnStyle, leftStyle, rightStyle color.Style,
	boundary, spaceSizeAfterLn, spaceSizeOnBoundary int,
) string {
	spaceSizeOnLn := len(rawLeftLn) // = len(rawRightLn)
	limit := boundary - (spaceSizeOnLn + spaceSizeAfterLn + spaceSizeOnBoundary)
	chl, chr := chopped(rawLeft, limit), chopped(rawRight, limit)

	w := &bytes.Buffer{}
	buf := bufio.NewWriter(w)

	longest := len(chl)
	if longest < len(chr) {
		longest = len(chr)
	}

	for i := 0; i < longest; i++ {
		leftText, rightText := get(i, chl), get(i, chr)

		var line string
		leftText = rpad(leftText, limit)
		if i == 0 {
			line = fmt.Sprintf(
				"%s%s%s%s%s%s%s\n",
				lnStyle.Sprint(rawLeftLn),
				strings.Repeat(" ", spaceSizeAfterLn),
				leftStyle.Sprint(leftText),
				strings.Repeat(" ", spaceSizeOnBoundary),
				lnStyle.Sprint(rawRightLn),
				strings.Repeat(" ", spaceSizeAfterLn),
				rightStyle.Sprint(rightText),
			)
		} else {
			line = fmt.Sprintf(
				"%s%s%s\n",
				leftStyle.Sprint(leftText),
				strings.Repeat(" ", spaceSizeOnBoundary),
				rightStyle.Sprint(rightText),
			)
		}
		buf.WriteString(line)
	}
	buf.Flush()
	return string(w.Bytes())
}

func terminalShape() (int, int, error) {
	var (
		out *os.File
		err error
		sz  struct {
			rows    uint16
			cols    uint16
			xpixels uint16
			ypixels uint16
		}
	)
	if runtime.GOOS == "openbsd" {
		out, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return 0, 0, err
		}
	} else {
		out, err = os.OpenFile("/dev/tty", os.O_WRONLY, 0)
		if err != nil {
			return 0, 0, err
		}
	}
	syscall.Syscall(
		syscall.SYS_IOCTL,
		out.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)),
	)
	return int(sz.cols), int(sz.rows), nil
}
