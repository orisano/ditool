// Copyright © 2017 Nao YONASHIRO <owan.orisano@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/cobra"
)

var (
	outputFile    string
	interfaceName string
)

var interfaceTemplate = template.Must(template.New("interface").Parse(`// GENERATED BY github.com/orisano/gendi
package {{.PkgName}}

type {{.InterfaceName}} interface {}
`))

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(interfaceName) == 0 {
			return errors.New("missing -i")
		}
		if len(outputFile) == 0 {
			outputFile = snaker.CamelToSnake(interfaceName) + ".go"
		}
		f, err := os.Create(outputFile)
		if err != nil {
			return errors.Wrap(err, "failed to create file")
		}
		defer f.Close()
		a, err := filepath.Abs(outputFile)
		if err != nil {
			return errors.Wrap(err, "failed to get abspath")
		}
		err = interfaceTemplate.Execute(f, struct{ InterfaceName, PkgName string }{
			InterfaceName: interfaceName,
			PkgName:       filepath.Base(filepath.Dir(a)),
		})
		if err != nil {
			return errors.Wrap(err, "template failed to execute")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&outputFile, "o", "", "generate file path")
	addCmd.Flags().StringVar(&interfaceName, "i", "", "interface name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
