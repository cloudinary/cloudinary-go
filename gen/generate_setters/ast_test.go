package main

import (
	"github.com/heimdalr/dag"
	"go/ast"
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

func TestAst_SearchTypesInAstFile(t *testing.T) {
	t.Skip("Not implemented yet")
	type args struct {
		file     *ast.File
		filename string
		graph    *dag.DAG
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
