package resize

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Crop resize
type cropGroup struct {
	width       interface{} `cld:"w"`
	height      interface{} `cld:"h"`
	aspectRatio interface{} `cld:"ar"`

	cropMode interface{} `cld:"c"`
	//Dimensions
	//Position
}

func (c cropGroup) Width(width int) cropGroup {
	c.width = width

	return c
}

func (c cropGroup) WidthPercent(width float64) cropGroup {
	c.width = width

	return c
}

func (c cropGroup) WidthExpr(width string) cropGroup {
	c.width = width

	return c
}

func (c cropGroup) Height(height int) cropGroup {
	c.height = height

	return c
}

func (c cropGroup) HeightPercent(height float64) cropGroup {
	c.height = height

	return c
}

func (c cropGroup) HeightExpr(height string) cropGroup {
	c.height = height

	return c
}

func (c cropGroup) AspectRatio(aspectRatio float64) cropGroup {
	c.aspectRatio = aspectRatio

	return c
}

func Crop() cropGroup {
	return cropGroup{cropMode: "crop"}
}

func Scale() cropGroup {
	return cropGroup{cropMode: "scale"}
}
func (c cropGroup) String() string {
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
