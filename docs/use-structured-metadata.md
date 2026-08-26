# Use structured metadata

## When to use

Attach typed, validated business fields (SKU, campaign, rights expiry) to assets, so
applications can filter and route on a stable schema instead of free-form tags.

Tags are free-form strings; structured metadata fields are declared once per product
environment with a type, and the server validates every write against that schema. Use tags
for loose grouping and metadata for anything your code branches on.

## Define a field, once per environment

```go
field, err := cld.Admin.AddMetadataField(ctx, metadata.Field{
	ExternalID: "sku",
	Label:      "SKU",
	Type:       metadata.StringFieldType,
	Mandatory:  false, // plain bool here, not *bool as in UploadParams
})
if err != nil {
	return fmt.Errorf("add field transport: %w", err)
}
if field.Error.Message != "" {
	return fmt.Errorf("add field rejected: %s", field.Error.Message)
}
```

Note `AddMetadataField` takes a `metadata.Field` **directly**, rather than a wrapper params
struct.

Field types live in the `metadata` package: `StringFieldType`, `IntegerFieldType`,
`DateFieldType`, `EnumFieldType`, `SetFieldType`. Enum and set fields need a `DataSource`
listing the permitted values.

Definitions are **permanent and per-environment**. Re-adding an existing `ExternalID`
returns `external id sku already exists` (verified) — with `err == nil`, so check
`result.Error.Message`. Treat field creation as a migration you run once, not something in
application startup.

## Write values

At upload time:

```go
result, err := cld.Upload.Upload(ctx, source, uploader.UploadParams{
	PublicID:  "examples/product-photo",
	Overwrite: api.Bool(true),
	Metadata:  api.Metadata{"sku": "SKU-00042"},
})
```

Or afterwards, on assets that already exist:

```go
updated, err := cld.Upload.UpdateMetadata(ctx, uploader.UpdateMetadataParams{
	Metadata:  api.CldAPIMap{"sku": "SKU-00042"},
	PublicIDs: []string{"examples/product-photo"},
})
```

Mind the type difference between the two paths — the compiler will catch it, but the names
give no hint:

| Where | Metadata type | IDs type |
|---|---|---|
| `uploader.UploadParams.Metadata` | `api.Metadata` (`map[string]interface{}`) | — |
| `uploader.UpdateMetadataParams.Metadata` | `api.CldAPIMap` (`map[string]string`) | `[]string` |

Both serialize to `key=value|key=value` on the wire, so values are sent as strings regardless
of the declared field type. The server does the type validation.

## Undefined keys reject the whole write

**An undefined metadata key is rejected rather than silently dropped**, so an upload never
half-succeeds with metadata missing. Verified: uploading with a key that has no field
definition rejects the entire upload —

```
Metadata External IDs do not exist: ["no_such_field_xyz"]
```

— and, as always, arrives with `err == nil` and an **empty** `PublicID`. No asset is created.
So a typo'd metadata key does not degrade gracefully; it loses the upload. Define fields
before referencing them, and check `result.Error.Message`.

This is the opposite of the parameter behaviour elsewhere in the SDK, where unset struct
fields are simply omitted from the request.

## Query by metadata

Metadata is indexed and searchable:

```go
result, err := cld.Admin.Search(ctx, search.Query{
	Expression: `metadata.sku="SKU-00042"`,
	MaxResults: 30,
})
```

Use a raw string literal or escape the quotes; the value needs them. As with all searches, a
non-matching or misspelled field name returns zero results rather than an error — see
[Search and manage assets](search-and-manage-assets.md#writing-expressions-that-match).

## Reading values back

Metadata comes back on the result as `api.Metadata`:

```go
details, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "examples/product-photo"})
// ... check err and details.Error.Message ...
fmt.Println(details.Metadata["sku"])
```

## Pattern: model output to reviewed metadata

A useful pipeline when you are deriving metadata from AI analysis:

1. Run analysis on the asset (`cld.Admin.Analyze`, or the
   [Analyze API](https://cloudinary.com/documentation/analyze_api_guide.md) — subscription
   required).
2. **Normalize against your schema** — map free-form output onto your allowed values, drop
   low-confidence results, apply business rules.
3. Write the result as structured metadata.
4. Search, route, and deliver based on that metadata.

Step 2 is where the value is, and it is not optional: enum and set fields reject values
outside their datasource, so unnormalized model output fails the write. Automate the mapping
where your rules are clear and route the rest to a person — not because human review is
inherently required, but because a rejected write costs you the whole upload.

## Troubleshooting

- `Metadata External IDs do not exist: [...]` — the field is not defined in this
  environment, or the external ID is misspelled. The upload failed entirely.
- `external id <id> already exists` — definitions are permanent; reuse the field rather than
  re-creating it.
- An enum/set write fails — the value is not in the datasource. Update it first with
  `cld.Admin.UpdateMetadataFieldDataSource`.
- A metadata search returns nothing — check the `metadata.` prefix and the quoting on the
  value.
- Values written but not visible on the asset — you checked only `err`. Check
  `result.Error.Message`; see [Handle errors](handle-errors.md).

## Related

- Runnable example: `examples/use-structured-metadata/main.go` (in the repository)
- [Search and manage assets](search-and-manage-assets.md)
- [Structured metadata guide](https://cloudinary.com/documentation/structured_metadata.md)
- [Metadata API reference](https://cloudinary.com/documentation/admin_api.md#structured_metadata_api)
