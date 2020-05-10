package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"

	"github.com/ynqa/diffy/pkg/diffy"
)

var (
	context  int
	noHeader bool
	style    string
	tabSize  int

	version = "unversioned"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diffy [flags] FILE1 FILE2",
		Short: "Print colored diff more readable",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate(args); err != nil {
				return err
			}
			return execute(args)
		},
		Version: version,
	}

	cmd.Flags().IntVarP(&context, "context", "c", 3, "number of context to print")
	cmd.Flags().BoolVar(&noHeader, "no-header", false, "no file name header")
	cmd.Flags().StringVarP(&style, "style", "s", "unified", "output style; one of unified|split")
	cmd.Flags().IntVar(&tabSize, "tab-size", 4, "tab stop spacing")
	return cmd
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func validate(args []string) error {
	if len(args) != 2 {
		return errors.New("length for args must be 2")
	}
	if !fileExists(args[0]) {
		return fmt.Errorf("file '%s' was not found", args[0])
	}
	if !fileExists(args[1]) {
		return fmt.Errorf("file '%s' was not found", args[1])
	}
	if context < 1 {
		return fmt.Errorf("got context '%d', but it must be positive", context)
	}
	if tabSize < 1 {
		return fmt.Errorf("got tabSize '%d', but it must be positive", tabSize)
	}
	if style != "unified" && style != "split" {
		return fmt.Errorf("got style '%s', but it must be one of unified|split", style)
	}
	return nil
}

func execute(args []string) error {
	org, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	new, err := ioutil.ReadFile(args[1])
	if err != nil {
		return err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(org)),
		B:        difflib.SplitLines(string(new)),
		FromFile: args[0],
		ToFile:   args[1],
		Context:  context,
	}
	opt := diffy.Option{
		NoHeader:         noHeader,
		TabSize:          tabSize,
		SeparatorSymbol:  "-",
		SeparatorWidth:   4,
		SpaceSizeAfterLn: 2,
	}
	if style == "unified" {
		return diffy.WriteUnifiedDiff(os.Stdin, diff, opt)
	} else if style == "split" {
		return diffy.WriteSplitDiff(os.Stdin, diff, opt)
	}
	return nil
}
