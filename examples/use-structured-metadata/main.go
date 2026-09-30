// Define a structured metadata field, write a value to an asset, and query by it.
//
// Structured metadata is typed and server-validated. An undefined key rejects the whole
// upload rather than being dropped, so define fields before writing values to them.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./use-structured-metadata
//
// In a real project the field definition is a one-off migration, not something your
// application runs at startup: definitions are permanent and per-environment, and
// re-adding one reports "external id already exists". This example tolerates that.
//
// Docs: docs/use-structured-metadata.md
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/metadata"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sourceURL = "https://res.cloudinary.com/demo/image/upload/sample.jpg"
	publicID  = "examples/product-photo"

	fieldExternalID = "sku"
	fieldLabel      = "SKU"
	skuValue        = "SKU-00042"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	cld, err := cloudinary.New()
	if err != nil {
		return err
	}

	ctx := context.Background()

	// 1. Define the field. AddMetadataField takes a metadata.Field directly, not a
	//    wrapper params struct.
	field, err := cld.Admin.AddMetadataField(ctx, metadata.Field{
		ExternalID: fieldExternalID,
		Label:      fieldLabel,
		Type:       metadata.StringFieldType,
		Mandatory:  false, // plain bool on metadata.Field
	})
	if err != nil {
		return fmt.Errorf("add field transport: %w", err)
	}
	switch {
	case field.Error.Message == "":
		fmt.Printf("defined field %q\n", field.ExternalID)
	case strings.Contains(field.Error.Message, "already exists"):
		// Expected on any run after the first — definitions are permanent.
		fmt.Printf("field %q already defined\n", fieldExternalID)
	default:
		return fmt.Errorf("add field rejected: %s", field.Error.Message)
	}

	// 2. Write a value at upload time. An undefined key here would fail the whole
	//    upload with "Metadata External IDs do not exist".
	uploaded, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: api.Bool(true),
		Metadata:  api.Metadata{fieldExternalID: skuValue},
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}
	fmt.Printf("uploaded %s with metadata %v\n", uploaded.PublicID, uploaded.Metadata)

	// 3. Or write to assets that already exist. Note the different types on this
	//    path: CldAPIMap rather than Metadata, []string rather than CldAPIArray.
	updated, err := cld.Upload.UpdateMetadata(ctx, uploader.UpdateMetadataParams{
		Metadata:  api.CldAPIMap{fieldExternalID: skuValue},
		PublicIDs: []string{publicID},
	})
	if err != nil {
		return fmt.Errorf("update metadata transport: %w", err)
	}
	// UpdateMetadataResult.Error is interface{}, not api.ErrorResp as on most result
	// types — compare against nil rather than reaching for .Message.
	if updated.Error != nil {
		return fmt.Errorf("update metadata rejected: %v", updated.Error)
	}

	// 4. Query by metadata. The value needs quoting inside the expression.
	found, err := cld.Admin.Search(ctx, search.Query{
		Expression: fmt.Sprintf(`metadata.%s=%q`, fieldExternalID, skuValue),
		MaxResults: 10,
	})
	if err != nil {
		return fmt.Errorf("search transport: %w", err)
	}
	if found.Error.Message != "" {
		return fmt.Errorf("search rejected: %s", found.Error.Message)
	}

	fmt.Printf("assets with %s=%s: %d\n", fieldExternalID, skuValue, found.TotalCount)
	for _, asset := range found.Assets {
		fmt.Printf("  %s\n", asset.PublicID)
	}
	if found.TotalCount == 0 {
		fmt.Println("(the search index lags writes by a moment; retry if this is empty)")
	}

	return nil
}
