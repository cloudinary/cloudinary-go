package main

import (
	"github.com/heimdalr/dag"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAst_GetPackageAndStructNamesByAstField(t *testing.T) {
	type args struct {
		currentPackage string
		field          *ast.Field
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Struct name field",
			args: args{
				currentPackage: "main",
				field: &ast.Field{
					Type: &ast.Ident{
						Name: "test",
					},
				},
			},
			want:  "main",
			want1: "test",
		},
		{
			name: "Struct from another package",
			args: args{
				currentPackage: "main",
				field: &ast.Field{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "not_main",
						},
						Sel: &ast.Ident{
							Name: "test1",
						},
					},
				},
			},
			want:  "not_main",
			want1: "test1",
		},
		{
			name: "Unknown field type",
			args: args{
				currentPackage: "main2",
				field: &ast.Field{
					Type: &ast.ArrayType{},
				},
			},
			want:  "",
			want1: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getPackageAndStructNamesByAstField(tt.args.currentPackage, tt.args.field)
			if got != tt.want {
				t.Errorf("getPackageAndStructNamesByAstField() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getPackageAndStructNamesByAstField() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_searchTypesInAstFile(t *testing.T) {
	type file struct {
		file     *ast.File
		filename string
	}

	type args struct {
		files []file
	}

	type expected struct {
		vertices []string
		edges    map[string][]string
	}

	type test struct {
		name     string
		args     args
		expected expected
	}

	var tests []test

	tests = append(tests, test{
		name: "Empty files causes empty graph",
		args: args{
			files: nil,
		},
		expected: expected{
			vertices: nil,
			edges:    nil,
		},
	})

	test2 := test{
		name: "One structure without embeding",
		args: args{
			files: []file{},
		},
		expected: expected{
			vertices: []string{"testdata.struct1"},
			edges:    nil,
		},
	}

	for _, filename := range []string{"./testdata/struct1.go"} {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		test2.args.files = append(test2.args.files, file{parsedFile, filename})
	}

	tests = append(tests, test2)

	test3 := test{
		name: "Two structures without embeding",
		args: args{
			files: []file{},
		},
		expected: expected{
			vertices: []string{"testdata.struct1", "testdata.struct2"},
			edges:    nil,
		},
	}

	for _, filename := range []string{"./testdata/struct1.go", "./testdata/struct2.go"} {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		test3.args.files = append(test3.args.files, file{parsedFile, filename})
	}

	tests = append(tests, test3)

	test4 := test{
		name: "Two structures, one from different package",
		args: args{
			files: []file{},
		},
		expected: expected{
			vertices: []string{"testdata.struct1", "anotherpackage.Struct3"},
			edges:    nil,
		},
	}

	for _, filename := range []string{"./testdata/struct1.go", "./testdata/anotherpackage/struct3.go"} {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		test4.args.files = append(test4.args.files, file{parsedFile, filename})
	}

	tests = append(tests, test4)

	test5 := test{
		name: "Cases file",
		args: args{
			files: []file{},
		},
		expected: expected{
			vertices: []string{"testdata.struct1", "testdata.struct2", "anotherpackage.Struct3", "testdata.case1", "testdata.case2", "testdata.case3"},
			edges: map[string][]string{
				"testdata.case1": {"testdata.struct1"},
				"testdata.case2": {"testdata.struct1", "testdata.struct2"},
				"testdata.case3": {"testdata.struct1", "anotherpackage.Struct3"},
			},
		},
	}

	for _, filename := range []string{"./testdata/cases.go", "./testdata/struct1.go", "./testdata/anotherpackage/struct3.go", "./testdata/struct2.go"} {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		test5.args.files = append(test5.args.files, file{parsedFile, filename})
	}

	tests = append(tests, test5)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := dag.NewDAG()
			for _, file := range tt.args.files {
				searchTypesInAstFile(file.file, file.filename, graph)
			}

			if graph.GetOrder() != len(tt.expected.vertices) {
				t.Errorf("Expected result graph to have order %d. %d given", len(tt.expected.vertices), graph.GetOrder())
			}

			for _, id := range tt.expected.vertices {
				if _, err := graph.GetVertex(id); err != nil {
					t.Errorf("Vertex %s is not found in the result graph", id)
				}
			}

			expectedSize := 0

			for src, dests := range tt.expected.edges {
				for _, dest := range dests {
					expectedSize++
					if isEdge, err := graph.IsEdge(src, dest); err != nil {
						t.Errorf("Unexpected error %v during edge check", err)
					} else if !isEdge {
						t.Errorf("Edge %s->%s is not found in the result graph", src, dest)
					}
				}
			}

			if graph.GetSize() != expectedSize {
				t.Errorf("Expected result graph to have size %d. %d given", expectedSize, graph.GetSize())
			}
		})
	}
}
