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
	"github.com/spf13/cobra"
)

var pkgName string
var srcTemplate = template.Must(template.New("src").Parse(`// GENERATED BY github.com/orisano/ditool
package {{.PkgName}}

import (
	"fmt"
)

type Dependencies struct {
	m map[interface{}]func(d *Dependencies) (interface{}, error)
}

func (d *Dependencies) Register(key interface{}, factory func(d *Dependencies) (interface{}, error)) {
	d.m[key] = factory
}

func (d *Dependencies) Get(key interface{}) (interface{}, error) {
	factory, ok := d.m[key]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return factory(d)
}
`))

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(pkgName) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to get working directory")
			}
			pkgName = filepath.Base(wd)
		}
		f, err := os.Create("dependencies.go")
		if err != nil {
			return errors.Wrap(err, "failed to create file 'dependencies.go'")
		}
		defer f.Close()
		err = srcTemplate.Execute(f, struct{ PkgName string }{
			PkgName: pkgName,
		})
		if err != nil {
			errors.Wrap(err, "template failed to execute")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&pkgName, "p", "", "dependencies.go package name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
