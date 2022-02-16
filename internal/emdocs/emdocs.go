/*
Copyright © 2020 Equinix Metal Developers <support@equinixmetal.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package emdocs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```sh\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```sh\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

// GenMarkdownEM creates custom markdown output.
func GenMarkdownEM(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Description\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}

func appendCmd(c *cobra.Command, f *os.File) error {
	if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
		return nil
	}
	err := GenMarkdownEM(c, f)
	return err
}

func allCmds(c *cobra.Command) []*cobra.Command {
	all := []*cobra.Command{}
	for _, c := range c.Commands() {
		all = append(all, c)
		all = append(all, allCmds(c)...)
	}
	return all
}

// GenMarkdownEMDocs is the the same as GenMarkdownEM, but
// with custom filePrepender and linkHandler.
func GenMarkdownEMDocs(cmd *cobra.Command, dir string) error {
	basename := strings.Split(cmd.Root().Use, " ")[0] + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, c := range allCmds(cmd) {
		if err != appendCmd(c, f) {
			return err
		}
	}

	if !cmd.DisableAutoGenTag {
		f.Write([]byte("###### Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "\n"))
	}

	return nil
}

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "emdocs [DESTINATION]",
		Short:                 "Generate command documentation",
		Long:                  "To generate documentation in the ./docs directory: emdocs ./docs",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			dest := args[0]
			return GenMarkdownEMDocs(cmd.Parent(), dest)
		},
	}
}
