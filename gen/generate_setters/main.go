package main

import (
	"fmt"
	"github.com/heimdalr/dag"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	curDir, _ := os.Getwd()
	graph := dag.NewDAG()

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".go" {
				fset := token.NewFileSet()
				file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}

				searchTypesInAstFile(file, path, graph)
			}
			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	for id, cldType := range graph.GetRoots() {
		setters := getSettersForType(id, graph)
		if len(setters) > 0 {
			t, ok := cldType.(*CldType)
			if !ok {
				panic("Graph node can not be converted to CldType")
			}

			fmt.Printf("Generating setters file for struct %s with %d setters\n", id, len(setters))
			tmpl := generateTemplateForSettersFile(setters)

			fileName := fmt.Sprintf("%s/%s/%s_setters.go", curDir, t.filePath, strings.ToLower(t.structName))
			file, _ := os.Create(fileName)
			if err = writeTemplateToFile(file, tmpl, t); err != nil {
				panic(err)
			}

			if err := file.Close(); err != nil {
				panic(err)
			}
			fmt.Printf("Done. File: %s\n", fileName)
		}
	}
}
