package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// Generate template for setters file
func generateTemplateForSettersFile(setters []CldSetter) string {
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
	tmpl := template.Must(template.New("test").Delims("<<", ">>").Parse(stringTemplate))
	if err := tmpl.Execute(w, templateData); err != nil {
		panic(err)
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	return buf.String()
}

// Write template to file filling it with type data
func writeTemplateToFile(file io.Writer, templateString string, t *CldType) error {
	tmpl := template.Must(template.New("setters").Parse(templateString))

	mixinData := struct {
		PackageName string
		StructName  string
		Receiver    string
	}{
		t.packageName,
		t.structName,
		strings.ToLower(t.structName)[0:1],
	}

	w := bufio.NewWriter(file)
	if err := tmpl.Execute(w, mixinData); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
