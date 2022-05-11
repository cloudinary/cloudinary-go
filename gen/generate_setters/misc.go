package main

import (
	"github.com/heimdalr/dag"
	"path/filepath"
	"strings"
)

// CldSetter is the setter.
type CldSetter struct {
	Type   string
	Suffix string
	Param  string
}

// CldType is the Cloudinary node type.
type CldType struct {
	packageName string
	structName  string
	id          string
	setters     []CldSetter
	filePath    string
}

//ID returns the ID of the node.
func (t *CldType) ID() string {
	return t.id
}

// Get all setters for type based on embedding graph
func getSettersForType(id string, graph *dag.DAG) []CldSetter {
	vertex, _ := graph.GetVertex(id)
	setters := vertex.(*CldType).setters

	if child, err := graph.GetChildren(id); err == nil {
		for id := range child {
			setters = append(setters, getSettersForType(id, graph)...)
		}
	}

	return setters
}

// Search in graph for the node by id. Creates new node if not found.
func findOrCreateGraphNode(id string, packageName string, typeName string, filename string, graph *dag.DAG) *CldType {
	var cldType *CldType
	if _, err := graph.GetVertex(id); err != nil {
		cldType = &CldType{
			id:          id,
			packageName: packageName,
			structName:  typeName,
			filePath:    filepath.Dir(filename),
		}

		if _, err := graph.AddVertex(cldType); err != nil {
			panic(err)
		}
	} else {
		t, _ := graph.GetVertex(id)
		cldType = t.(*CldType)
	}

	return cldType
}

// Parse annotation and get CldSetters from it
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
