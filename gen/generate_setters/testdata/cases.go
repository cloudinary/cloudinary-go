package testdata

import "github.com/cloudinary/cloudinary-go/v2/gen/generate_setters/testdata/anotherpackage"

type case1 struct {
	struct1
}

type case2 struct {
	struct1
	struct2
}

type case3 struct {
	struct1
	anotherpackage.Struct3
}
