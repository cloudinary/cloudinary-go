package main

import (
	"github.com/heimdalr/dag"
	"reflect"
	"testing"
)

func Test_GetSettersForType(t *testing.T) {
	graph := dag.NewDAG()

	type1 := &CldType{
		packageName: "testPackage",
		structName:  "testStruct",
		id:          "testpackage_teststruct",
		setters: []CldSetter{
			{
				Type:   "int",
				Suffix: "suffix",
				Param:  "param1",
			},
		},
		filePath: "testFilePath",
	}

	graph.AddVertex(type1)

	type2 := &CldType{
		packageName: "testPackage",
		structName:  "testStruct2",
		id:          "testpackage_teststruct2",
		setters: []CldSetter{
			{
				Type:   "string",
				Suffix: "suffix2",
				Param:  "param2",
			},
		},
	}

	graph.AddVertex(type2)
	graph.AddEdge(type1.ID(), type2.ID())

	type3 := &CldType{
		packageName: "testPackage",
		structName:  "testStruct3",
		id:          "testpackage_teststruct3",
		setters: []CldSetter{
			{
				Type:   "string",
				Suffix: "suffix3",
				Param:  "param3",
			},
		},
	}

	graph.AddVertex(type3)

	type args struct {
		id    string
		graph *dag.DAG
	}
	tests := []struct {
		name string
		args args
		want []CldSetter
	}{
		{
			name: "testpackage_teststruct",
			args: args{
				id:    "testpackage_teststruct",
				graph: graph,
			},
			want: []CldSetter{},
		},
		{
			name: "testpackage_teststruct2",
			args: args{
				id:    "testpackage_teststruct2",
				graph: graph,
			},
			want: []CldSetter{},
		},
		{
			name: "testpackage_teststruct3",
			args: args{
				id:    "testpackage_teststruct3",
				graph: graph,
			},
			want: []CldSetter{},
		},
	}

	tests[0].want = append(tests[0].want, type1.setters...)
	tests[0].want = append(tests[0].want, type2.setters...)

	tests[1].want = append(tests[1].want, type2.setters...)

	tests[2].want = append(tests[2].want, type3.setters...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSettersForType(tt.args.id, tt.args.graph); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSettersForType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindOrCreateGraphNode(t *testing.T) {
	type args struct {
		id          string
		packageName string
		typeName    string
		filename    string
		graph       *dag.DAG
	}

	type test struct {
		name string
		args args
		want *CldType
	}

	var tests []test

	tests = append(tests, test{
		name: "Create new element in graph",
		args: args{
			id:          "test",
			packageName: "main",
			typeName:    "type",
			filename:    "test1/test2/filename",
			graph:       dag.NewDAG(),
		},
		want: &CldType{
			packageName: "main",
			structName:  "type",
			id:          "test",
			setters:     nil,
			filePath:    "test1/test2",
		},
	},
	)

	type1 := &CldType{
		packageName: "main2",
		structName:  "type2",
		id:          "test2",
		setters:     nil,
		filePath:    "test1/test2",
	}

	graph := dag.NewDAG()
	graph.AddVertex(type1)

	tests = append(tests, test{
		name: "Get existing element from graph",
		args: args{
			id:          "test2",
			packageName: "main",
			typeName:    "type",
			filename:    "test1/test2/test3/filename",
			graph:       graph,
		},
		want: type1,
	},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOrCreateGraphNode(tt.args.id, tt.args.packageName, tt.args.typeName, tt.args.filename, tt.args.graph); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findOrCreateGraphNode() = %v, want %v", got, tt.want)
			}

			if _, err := tt.args.graph.GetVertex(tt.args.id); err != nil {
				t.Errorf("Graph should contains a vertex. Error while getting it: %v", err)
			}
		})
	}
}

func Test_GetSettersFromAnnotation(t *testing.T) {
	type args struct {
		annotation string
		paramName  string
	}
	tests := []struct {
		name string
		args args
		want []CldSetter
	}{
		{
			name: "Single element without suffixes",
			args: args{
				annotation: "int",
				paramName:  "param",
			},
			want: []CldSetter{
				{
					Type:   "int",
					Suffix: "",
					Param:  "param",
				},
			},
		},
		{
			name: "Multiple elements without suffixes",
			args: args{
				annotation: "int,string",
				paramName:  "param",
			},
			want: []CldSetter{
				{
					Type:   "int",
					Suffix: "",
					Param:  "param",
				},
				{
					Type:   "string",
					Suffix: "",
					Param:  "param",
				},
			},
		},
		{
			name: "Single element with suffix",
			args: args{
				annotation: "float32:Percent",
				paramName:  "param",
			},
			want: []CldSetter{
				{
					Type:   "float32",
					Suffix: "Percent",
					Param:  "param",
				},
			},
		},
		{
			name: "Multiple elements with suffixes",
			args: args{
				annotation: "float32:Percent,string:Expr",
				paramName:  "param",
			},
			want: []CldSetter{
				{
					Type:   "float32",
					Suffix: "Percent",
					Param:  "param",
				},
				{
					Type:   "string",
					Suffix: "Expr",
					Param:  "param",
				},
			},
		},
		{
			name: "Multiple mixed elements",
			args: args{
				annotation: "int,float32:Percent,string",
				paramName:  "param",
			},
			want: []CldSetter{
				{
					Type:   "int",
					Suffix: "",
					Param:  "param",
				},
				{
					Type:   "float32",
					Suffix: "Percent",
					Param:  "param",
				},
				{
					Type:   "string",
					Suffix: "",
					Param:  "param",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSettersFromAnnotation(tt.args.annotation, tt.args.paramName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSettersFromAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Returns true on exists",
			args: args{
				s:   []string{"string1", "string2"},
				str: "string1",
			},
			want: true,
		},
		{
			name: "Returns false on not exists",
			args: args{
				s:   []string{"string1", "string2"},
				str: "string3",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
