package resize

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Crop resize
type CropGroup struct {
	cropMode interface{} `cld:"c"`
	Dimensions
}

func Crop() *CropGroup {
	return &CropGroup{cropMode: "crop"}
}

func Scale() *CropGroup {
	return &CropGroup{cropMode: "scale"}
}

func (c *CropGroup) String() string {
	v := reflect.ValueOf(*c)
	result := stringifyReflectedValue(v)
	sort.Strings(result)

	return strings.Join(result, ",")
}

func stringifyReflectedValue(v reflect.Value) []string {
	var res []string

	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if isEmptyValue(v.Field(i)) {
			continue
		}

		if v.Field(i).Kind() == reflect.Struct {
			res = append(res, stringifyReflectedValue(v.Field(i))...)
		} else {
			res = append(res, fmt.Sprintf("%s_%v", typeOfS.Field(i).Tag.Get("cld"), v.Field(i).Elem()))
		}
	}

	return res
}
