package main

import (
	"fmt"
	"github.com/heimdalr/dag"
	"go/ast"
	"reflect"
)

// Parse reflected field. Returns currentPackage if field name does not contain other package name.
func getPackageAndStructNamesByAstField(currentPackage string, field *ast.Field) (string, string) {
	var embedStructName string
	var embedPackageName string

	if embedStruct, ok := field.Type.(*ast.Ident); ok {
		embedStructName = embedStruct.Name
		embedPackageName = currentPackage
	}

	if selector, ok := field.Type.(*ast.SelectorExpr); ok {
		if embedPackage, ok := selector.X.(*ast.Ident); ok {
			embedPackageName = embedPackage.Name
		}

		embedStructName = selector.Sel.Name
	}

	return embedPackageName, embedStructName
}

// Parses a given AST file and adds all found types to the given graph
func searchTypesInAstFile(file *ast.File, filename string, graph *dag.DAG) {
	var currentTypeName string
	var currentPackage string

	ast.Inspect(file, func(x ast.Node) bool {
		if pkg, ok := x.(*ast.File); ok {
			currentPackage = pkg.Name.String()
		}

		if ss, ok := x.(*ast.TypeSpec); ok {
			currentTypeName = ss.Name.String()
		}

		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		structVertexId := fmt.Sprintf("%s.%s", currentPackage, currentTypeName)
		cldType := findOrCreateGraphNode(structVertexId, currentPackage, currentTypeName, filename, graph)

		for _, field := range s.Fields.List {
			if field.Names == nil {
				// Embed struct
				embedPackageName, embedStructName := getPackageAndStructNamesByAstField(currentPackage, field)
				embedVertexId := fmt.Sprintf("%s.%s", embedPackageName, embedStructName)
				findOrCreateGraphNode(embedVertexId, embedPackageName, embedStructName, filename, graph)
				if err := graph.AddEdge(structVertexId, embedVertexId); err != nil {
					panic(err)
				}
			} else if field.Names != nil && field.Tag != nil {
				// Just a field
				tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				if tag.Get("setters") != "" {
					cldType.setters = append(cldType.setters, getSettersFromAnnotation(tag.Get("setters"), field.Names[0].Name)...)
				}
			}
		}

		return false
	})
}
