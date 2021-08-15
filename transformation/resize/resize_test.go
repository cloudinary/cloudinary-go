package resize_test

import (
	"github.com/cloudinary/cloudinary-go/transformation/resize"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResize_Scale(t *testing.T) {
	s := resize.Scale().WidthPercent(0.5).Height(200)

	assert.Equal(t, "c_scale,h_200,w_0.5", s.String())
}

func TestResize_Fit(t *testing.T) {
	f := resize.Fit().AspectRatio(1).WidthPercent(1.0)

	assert.Equal(t, "ar_1,c_fit,w_1", f.String())
}

func TestResize_Crop(t *testing.T) {
	c := resize.Crop().HeightPercent(0.5).AspectRatio(2).WidthExpr("iw_div_2").X(50).Y(100)

	assert.Equal(t, "ar_2,c_crop,h_0.5,w_iw_div_2,x_50,y_100", c.String())
}


func TestResize_Thumbnail(t *testing.T) {
	th := resize.Thumbnail().Height(200).AspectRatioExpr("16:9").YPercent(0.5)

	assert.Equal(t, "ar_16:9,c_thumb,h_200,y_0.5", th.String())
}
