package asset_test

import (
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/asset"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"log"
	"testing"
)

func TestAsset_String(t *testing.T) {
	i, err := asset.Image(cldtest.PublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	i.DeliveryType = api.Authenticated
	i.Config.URL.SignURL = true
	log.Println(i.String())

	v, err := asset.Video(cldtest.VideoPublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(v.String())

	f, err := asset.File("sample_file", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(f.String())

	m, err := asset.Media("test/" + cldtest.PublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	m.Transformation = "c_scale,w_500"
	log.Println(m.String())
}
