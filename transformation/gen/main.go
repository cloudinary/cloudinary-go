package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/heimdalr/dag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

func main() {
	curDir, _ := os.Getwd()

	d := dag.NewDAG()

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".go" {
				searchTypesInFile(path, d)
			}
			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	for id, cldType := range d.GetRoots() {
		setters := getSettersForType(id, d)
		if len(setters) > 0 {
			fmt.Printf("Generating setters file for struct %s with %d setters\n", id, len(setters))
			template := getSettersTemplate(setters)
			fileName := generateMixinFile(curDir, template, cldType.(*CldType))
			fmt.Printf("Done. File: %s\n", fileName)
		}
	}
}

type CldSetter struct {
	Type   string
	Suffix string
	Param  string
}

type CldType struct {
	packageName string
	structName  string
	id          string
	setters     []CldSetter
	filePath    string
}

func (t *CldType) ID() string {
	return t.id
}

func getSettersForType(id string, d *dag.DAG) []CldSetter {
	vertex, _ := d.GetVertex(id)
	setters := vertex.(*CldType).setters

	if child, err := d.GetChildren(id); err == nil {
		for id, _ := range child {
			setters = append(setters, getSettersForType(id, d)...)
		}
	}

	return setters
}

func searchTypesInFile(filename string, d *dag.DAG) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var currentTypeName string
	var currentPackage string

	ast.Inspect(file, func(x ast.Node) bool {
		pkg, ok := x.(*ast.File)
		if ok {
			currentPackage = pkg.Name.String()
		}

		ss, ok := x.(*ast.TypeSpec)

		if ok {
			currentTypeName = ss.Name.String()
		}

		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		id := fmt.Sprintf("%s.%s", currentPackage, currentTypeName)
		var cldType *CldType

		if _, err := d.GetVertex(id); err != nil {
			cldType = &CldType{
				id:          id,
				packageName: currentPackage,
				structName:  currentTypeName,
				filePath:    filepath.Dir(filename),
			}

			d.AddVertex(cldType)
		} else {
			t, _ := d.GetVertex(id)
			cldType = t.(*CldType)
		}

		for _, field := range s.Fields.List {
			if field.Names == nil {
				var embedStructName string
				var embedPackageName string

				if embedStruct, ok := field.Type.(*ast.Ident); ok {
					embedStructName = embedStruct.Name
					embedPackageName = currentPackage
				}

				if selector, ok := field.Type.(*ast.SelectorExpr); ok {
					if embedPackage, ok := selector.X.(*ast.Ident); ok {
						embedPackageName = embedPackage.Name
					} else {
						continue
					}

					embedStructName = selector.Sel.Name
				}

				embedId := fmt.Sprintf("%s.%s", embedPackageName, embedStructName)

				if _, err := d.GetVertex(embedId); err != nil {
					t := CldType{
						id:          embedId,
						packageName: embedPackageName,
						structName:  embedStructName,
						filePath:    filepath.Dir(filename),
					}

					d.AddVertex(&t)
				}

				if err := d.AddEdge(id, embedId); err != nil {
					panic(err)
				}
			} else {
				if field.Tag != nil {
					sz := len(field.Tag.Value)
					tag := reflect.StructTag(field.Tag.Value[1 : sz-1])
					fieldName := field.Names
					if fieldName != nil && tag.Get("setters") != "" {
						cldType.setters = append(cldType.setters, getSettersFromAnnotation(tag.Get("setters"), fieldName[0].Name)...)
					}
				}
			}
		}

		return false
	})
}

func getSettersFromAnnotation(annotation string, paramName string) []CldSetter {
	var res []CldSetter
	types := strings.Split(annotation, ",")

	for _, t := range types {
		s := strings.Split(t, ":")
		if len(s) > 1 {
			res = append(res, CldSetter{s[0], s[1], paramName})
		} else {
			res = append(res, CldSetter{t, "", paramName})
		}
	}

	return res
}

func getSettersTemplate(setters []CldSetter) string {
	stringTemplate := "package {{ .PackageName}}\n"
	functionTemplate := `
func ({{ .Receiver}} *{{ .StructName}}) <<index .FuncName %d>>(<<index .ParamName %d>> <<index .Type %d>>) *{{ .StructName}} {
	{{ .Receiver}}.<<index .ParamName %d>> = <<index .ParamName %d>>

	return {{ .Receiver}}
}
`
	templateData := struct {
		FuncName  map[int]string
		ParamName map[int]string
		Type      map[int]string
	}{
		map[int]string{},
		map[int]string{},
		map[int]string{},
	}

	for i, t := range setters {
		stringTemplate += fmt.Sprintf(functionTemplate, i, i, i, i, i)
		templateData.FuncName[i] = strings.Title(t.Param) + t.Suffix
		templateData.Type[i] = t.Type
		templateData.ParamName[i] = t.Param
	}

	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	templ := template.Must(template.New("test").Delims("<<", ">>").Parse(stringTemplate))
	templ.Execute(w, templateData)
	w.Flush()

	return buf.String()
}

func generateMixinFile(curDir string, templateString string, t *CldType) string {
	templ := template.Must(template.New("test").Parse(templateString))

	mixinData := struct {
		PackageName string
		StructName  string
		Receiver    string
	}{
		t.packageName,
		t.structName,
		strings.ToLower(t.structName)[0:1],
	}

	fileName := fmt.Sprintf("%s/%s/%s_setters.go", curDir, t.filePath, strings.ToLower(t.structName))
	f, _ := os.Create(fileName)
	w := bufio.NewWriter(f)
	templ.Execute(w, mixinData)
	w.Flush()
	f.Close()

	return fileName
}
