package resize

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Crop resize
type CropGroup struct {
	cropMode   interface{} `cld:"c"`
	Dimensions `mixin:"true"`
}

func Crop() *CropGroup {
	return &CropGroup{cropMode: "crop"}
}

func Scale() *CropGroup {
	return &CropGroup{cropMode: "scale"}
}

func (c *CropGroup) String() string {
	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	var res []string
	for i := 0; i < v.NumField(); i++ {
		if isEmptyValue(v.Field(i)) {
			continue
		}

		res = append(res, fmt.Sprintf("%s_%v", typeOfS.Field(i).Tag.Get("cld"), v.Field(i).Elem()))
	}

	sort.Strings(res)

	return strings.Join(res, ",")
}
