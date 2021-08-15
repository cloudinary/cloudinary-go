package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestTemplate_GenerateTemplateForSettersFile(t *testing.T) {
	type args struct {
		setters []CldSetter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "One setter without suffix",
			args: args{
				setters: []CldSetter{{
					Type:   "int",
					Suffix: "",
					Param:  "param1",
				}},
			},
			want: `package {{ .PackageName}}

func ({{ .Receiver}} *{{ .StructName}}) Param1(param1 int) *{{ .StructName}} {
	{{ .Receiver}}.param1 = param1

	return {{ .Receiver}}
}
`,
		},
		{
			name: "One setter with suffix",
			args: args{
				setters: []CldSetter{{
					Type:   "int",
					Suffix: "Suffix",
					Param:  "param1",
				}},
			},
			want: `package {{ .PackageName}}

func ({{ .Receiver}} *{{ .StructName}}) Param1Suffix(param1 int) *{{ .StructName}} {
	{{ .Receiver}}.param1 = param1

	return {{ .Receiver}}
}
`,
		},
		{
			name: "Multiple mixed setters",
			args: args{
				setters: []CldSetter{
					{
						Type:   "int",
						Suffix: "",
						Param:  "param1",
					},
					{
						Type:   "string",
						Suffix: "Suffix",
						Param:  "param1",
					},
					{
						Type:   "float32",
						Suffix: "",
						Param:  "param2",
					},
				},
			},
			want: `package {{ .PackageName}}

func ({{ .Receiver}} *{{ .StructName}}) Param1(param1 int) *{{ .StructName}} {
	{{ .Receiver}}.param1 = param1

	return {{ .Receiver}}
}

func ({{ .Receiver}} *{{ .StructName}}) Param1Suffix(param1 string) *{{ .StructName}} {
	{{ .Receiver}}.param1 = param1

	return {{ .Receiver}}
}

func ({{ .Receiver}} *{{ .StructName}}) Param2(param2 float32) *{{ .StructName}} {
	{{ .Receiver}}.param2 = param2

	return {{ .Receiver}}
}
`,
		},
		{
			name: "Nil slice input",
			args: args{
				setters: nil,
			},
			want: "package {{ .PackageName}}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateTemplateForSettersFile(tt.args.setters); got != tt.want {
				t.Errorf("generateTemplateForSettersFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplate_WriteTemplateToFile(t *testing.T) {
	cldType := CldType{
		packageName: "package",
		structName:  "struct",
		id:          "id",
		setters: []CldSetter{
			{
				Type:   "int",
				Suffix: "suffix",
				Param:  "param",
			},
		},
		filePath: "test1/test2/test3",
	}

	type args struct {
		buffer         bytes.Buffer
		templateString string
		t              *CldType
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantString string
	}{
		{
			name: "Empty template = empty buffer",
			args: args{
				buffer:         bytes.Buffer{},
				templateString: "",
				t:              &cldType,
			},
			wantErr:    false,
			wantString: "",
		},
		{
			name: "Receiver in template subsitutes with small first letter of the struct name",
			args: args{
				buffer:         bytes.Buffer{},
				templateString: "{{ .Receiver}}",
				t:              &cldType,
			},
			wantErr:    false,
			wantString: "s",
		},
		{
			name: "Package name in template subsitutes with package name from struct",
			args: args{
				buffer:         bytes.Buffer{},
				templateString: "{{ .PackageName}}",
				t:              &cldType,
			},
			wantErr:    false,
			wantString: "package",
		},
		{
			name: "Struct name in template subsitutes with struct name",
			args: args{
				buffer:         bytes.Buffer{},
				templateString: "{{ .StructName}}",
				t:              &cldType,
			},
			wantErr:    false,
			wantString: "struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := bufio.NewWriter(&tt.args.buffer)
			if err := writeTemplateToFile(writer, tt.args.templateString, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("writeTemplateToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.buffer.String() != tt.wantString {
				t.Errorf("writeTemplateToFile should write %v to given buffer. %v given", tt.wantString, tt.args.buffer.String())
			}
		})
	}
}
