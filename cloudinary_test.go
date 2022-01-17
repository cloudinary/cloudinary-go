package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"log"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var c, _ = New()
var ctx = context.Background()

func TestCloudinary_CreateInstance(t *testing.T) {
	c, _ := New()

	if c.Config.Cloud.CloudName == "" {
		t.Error("Please set up CLOUDINARY_URL environment variable to run the test.")
	}

	c, _ = NewFromURL("cloudinary://key:secret@test123")

	if c.Config.Cloud.CloudName != "test123" {
		t.Error("Failed creating Cloudinary instance from Cloudinary URL.")
	}

	c, err := NewFromURL("")
	if err == nil || err.Error() != "must provide CLOUDINARY_URL" {
		t.Error("Error expected, got: ", err)
	}

	c, _ = NewFromParams("test123", "key", "secret")

	if c.Config.Cloud.CloudName != "test123" {
		t.Error("Failed creating Cloudinary instance from parameters.")
	}

	c, _ = NewFromOAuthToken("test123", "MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4")

	if c.Config.Cloud.CloudName != "test123" {
		t.Error("Failed creating Cloudinary instance from OAuth token.")
	}
}

func TestCloudinary_Upload(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:       "test_image",
		Eager:          "w_500,h_500",
		UniqueFilename: false,
		UseFilename:    true,
		Overwrite:      true,
	}

	resp, err := c.Upload.Upload(ctx, "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png", params)

	if err != nil {
		t.Error("Uploader failed: ", err)
	}

	if resp == nil || resp.PublicID != "test_image" {
		t.Error("Uploader failed: ", resp)
	}
}

func TestCloudinary_Admin(t *testing.T) {
	resp, err := c.Admin.Ping(ctx)

	if err != nil {
		t.Error("Admin API failed: ", err)
	}

	if resp == nil || resp.Status != "ok" {
		t.Error("Admin API failed: ", resp)
	}
}

func TestCloudinary_Asset(t *testing.T) {
	c.Config.URL.Secure = true
	c.Config.URL.SecureCDNSubDomain = true

	i, err := c.Image(cldtest.PublicID)
	if err != nil {
		t.Fatal(err)
	}
	i.DeliveryType = "fetch"
	i.Version = 123
	log.Println(i.String())

	v, err := c.Video(cldtest.VideoPublicID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(v.String())

	f, err := c.File("sample_file")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(f.String())

	m, err := c.Media(cldtest.PublicID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(m.String())
}
