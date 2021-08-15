package test_data

import "github.com/cloudinary/cloudinary-go/gen/generate_setters/test_data/another_package"

type case1 struct {
	struct1
}

type case2 struct {
	struct1
	struct2
}

type case3 struct {
	struct1
	another_package.Struct3
}
