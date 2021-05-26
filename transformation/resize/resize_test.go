package resize_test

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/transformation/resize"
	"testing"
)

func TestResize_Crop(t *testing.T) {
	c := resize.Crop().HeightPercent(0.5).AspectRatio(2).WidthExpr("iw_div_2")
	fmt.Printf("%v\n", c)

	s := resize.Scale().WidthPercent(0.5).Height(200)

	fmt.Printf("%v\n", s)
}
