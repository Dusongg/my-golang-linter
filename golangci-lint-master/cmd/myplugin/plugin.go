// file: cmd/myplugin/plugin.go
package myplugin

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// MyLinter implements a simple linter that looks for "TODO" comments
var MyLinter = &analysis.Analyzer{
	Name: "myplugin",
	Doc:  "checks for TODO comments in the code",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			for _, commentGroup := range file.Comments {
				for _, comment := range commentGroup.List {
					if strings.Contains(comment.Text, "TODO") {
						pass.Reportf(comment.Pos(), "found TODO comment: %s", comment.Text)
					}
				}
			}
		}
		return nil, nil
	},
}
