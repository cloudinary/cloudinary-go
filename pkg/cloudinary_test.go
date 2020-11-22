package cloudinary

import (
	"cloudinary-labs/cloudinary-go/pkg/api/uploader"
	"context"
	"testing"
)

var c, _ = Create()
var ctx = context.Background()

func TestCloudinary_CreateInstance(t *testing.T) {
	c, _ := Create()

	if c.Config.Account.CloudName == "" {
		t.Error("Something went wrong with Cloudinary instance")
	}

	c, _ = CreateFromUrl("cloudinary://key:secret@test123")

	if c.Config.Account.CloudName != "test123" {
		t.Error("Something went wrong with Cloudinary instance")
	}

	c, _ = CreateFromParams("test123", "key", "secret")

	if c.Config.Account.CloudName != "test123" {
		t.Error("Something went wrong with Cloudinary instance")
	}
}

func TestCloudinary_Upload(t *testing.T) {
	params := uploader.UploadParams{
		PublicID: "test_image",
	}

	resp, err := c.Upload.Upload("https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png", params)

	if err != nil {
		t.Error("Something went wrong with the uploader", err)
	}

	if resp == nil || resp.PublicID != "test_image" {
		t.Error("Something went wrong with the uploader", resp)
	}
}

func TestCloudinary_Admin(t *testing.T) {
	resp, err := c.Admin.Ping(ctx)

	if err != nil {
		t.Error("Something went wrong with admin api", err)
	}

	if resp == nil || resp.Status != "ok" {
		t.Error("Something went wrong with admin api", resp)
	}
}
