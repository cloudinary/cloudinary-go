package cloudinary

import (
	"context"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var c, _ = New()
var ctx = context.Background()

func TestCloudinary_CreateInstance(t *testing.T) {
	c, _ := New()

	if c.Config.Account.CloudName == "" {
		t.Error("Please set up CLOUDINARY_URL environment variable to run the test.")
	}

	c, _ = CreateFromUrl("cloudinary://key:secret@test123")

	if c.Config.Account.CloudName != "test123" {
		t.Error("Failed creating Cloudinary instance from Cloudinary URL.")
	}

	c, _ = CreateFromParams("test123", "key", "secret")

	if c.Config.Account.CloudName != "test123" {
		t.Error("Failed creating Cloudinary instance from parameters.")
	}
}

func TestCloudinary_Upload(t *testing.T) {
	params := uploader.UploadParams{
		PublicID: "test_image",
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
