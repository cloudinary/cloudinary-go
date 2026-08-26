# Search and manage assets

## When to use

Find assets by indexed fields, read or update attributes, and administer your media library
from the server. These use the Admin and Search APIs, which are **rate-limited** — treat them
as management operations, not a per-request database.

## Search

Expressions use Cloudinary's search syntax — fields, operators, ranges, and boolean
combinations are documented in the
[search expression reference](https://cloudinary.com/documentation/search_expressions.md).
The Go SDK takes the expression as a string in a `search.Query` struct; there is no fluent
builder.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
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

	query := search.Query{
		Expression: "resource_type:image AND created_at>1d",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 30,
	}

	result, err := cld.Admin.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("search transport: %w", err)
	}
	if result.Error.Message != "" {
		return fmt.Errorf("search rejected: %s", result.Error.Message)
	}

	fmt.Printf("total matches: %d\n", result.TotalCount)
	for _, asset := range result.Assets {
		fmt.Println(asset.AssetID, asset.PublicID, asset.Bytes, asset.CreatedAt)
	}

	// Pagination: feed NextCursor back until it comes back empty.
	if result.NextCursor != "" {
		query.NextCursor = result.NextCursor
		page2, err := cld.Admin.Search(ctx, query)
		if err != nil {
			return err
		}
		if page2.Error.Message != "" {
			return fmt.Errorf("search rejected: %s", page2.Error.Message)
		}
		fmt.Printf("second page: %d asset(s)\n", len(page2.Assets))
	}

	return nil
}
```

## Writing expressions that match

A syntactically valid expression that matches nothing returns `TotalCount: 0` with an empty
`Error.Message` — matching nothing is a valid answer, not an error. A few behaviours are
worth knowing when a query returns less than you expected:

| Expression | Result | Why |
|---|---|---|
| `resource_type:image` | matches | a plain field query works |
| `folder:examples` | 0 matches even with assets there | On a dynamic-folder environment, use `asset_folder:` or filter on `public_id` |
| `*` | rejected | `Query Error (at position 1)` — a bare wildcard is not a valid expression |
| `*ample` | rejected | Leading wildcards are not supported |

`folder:` is the one to watch, since it is valid syntax on every environment but only matches
on fixed-folder ones. If a search you expect to match returns zero, check the field name
against the [expression reference](https://cloudinary.com/documentation/search_expressions.md)
before concluding the assets are absent.

Note also that **the search index lags writes** by a short interval. For read-after-write,
use `cld.Admin.AssetByAssetID` instead of searching.

## List assets without a query

When you want everything of a type rather than a search expression, `cld.Admin.Assets`
paginates the same way — feed `NextCursor` back until it comes back empty:

```go
nextCursor := ""
for {
	page, err := cld.Admin.Assets(ctx, admin.AssetsParams{
		AssetType:  api.Image,
		MaxResults: 100,
		NextCursor: nextCursor,
	})
	if err != nil {
		return fmt.Errorf("list transport: %w", err)
	}
	if page.Error.Message != "" {
		return fmt.Errorf("list rejected: %s", page.Error.Message)
	}

	for _, asset := range page.Assets {
		fmt.Println(asset.PublicID, asset.SecureURL)
	}

	if page.NextCursor == "" {
		break
	}
	nextCursor = page.NextCursor
}
```

`MaxResults` caps at 500 per page. Bound the loop when you only need a sample, and remember
`AssetType` defaults to image — pass `api.Video` or `api.File` for the others. For anything
filtered, prefer `Search`: listing every asset to filter client-side burns rate limit.

## Read and update a single asset

```go
// By the immutable handle — survives renames and moves.
details, err := cld.Admin.AssetByAssetID(ctx, admin.AssetByAssetIDParams{AssetID: storedAssetID})
if err != nil {
	return err
}
if details.Error.Message != "" {
	return fmt.Errorf("lookup rejected: %s", details.Error.Message)
}

// Updates need the PUBLIC id — UpdateAsset has no asset-id variant.
updated, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
	PublicID: details.PublicID,
	Tags:     api.CldAPIArray{"featured"},
})
```

**Store the `AssetID` — it survives renames and moves — and resolve it to a `PublicID` when
you need to act on the asset.** `AssetByAssetID` is the asset-ID read entry point;
`UpdateAsset`, `Destroy`, `Rename`, and URL building take the public ID, so the lookup above
is the step that connects the two. Note that `cld.Admin.AssetsByIDs` takes **public** IDs.

The related-assets calls also accept asset IDs directly
(`AddRelatedAssetsByAssetIDs`, `DeleteRelatedAssetsByAssetIDs`).

## Remember the asset-type default

`admin.AssetParams`, `admin.UpdateAssetParams`, and friends default to **image**. Looking up
a video without setting `AssetType` returns `Resource not found`:

```go
cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "docs/vid", AssetType: api.Video})
```

The uploader calls the same concept `ResourceType`. Set the type explicitly on both sides
whenever you work with video or raw files.

## Deletion — destructive, no undo without backups

```go
cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: "examples/one-asset"})
cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: api.CldAPIArray{"a", "b"}})
cld.Admin.DeleteAssetsByPrefix(ctx, admin.DeleteAssetsByPrefixParams{Prefix: api.CldAPIArray{"examples/"}})
cld.Admin.DeleteAllAssets(ctx, admin.DeleteAllAssetsParams{}) // everything of that type
```

Prefer explicit ID lists: a prefix that matches more than you intended deletes exactly what
it matched and reports success, so `DeleteAssetsByPrefix` and `DeleteAllAssets` are worth
dry-running as a search first. Enable backups on the product environment to keep
`cld.Admin.RestoreAssets` available.

## Folders

```go
created, err := cld.Admin.CreateFolder(ctx, admin.CreateFolderParams{Folder: "examples/archive"})
// created.Success, created.Path, created.Name

renamed, err := cld.Admin.RenameFolder(ctx, admin.RenameFolderParams{
	FromPath: "examples/archive",
	ToPath:   "examples/archive-2024",
})
// renamed.From.Path, renamed.To.Path
```

The two results have different shapes: `CreateFolderResult` reports `Success`, while
`RenameFolderResult` returns `From` and `To` folder records instead. Read the paths off
`From`/`To` to confirm a rename, and `Error.Message` for rejections on both.

On dynamic-folder environments the asset folder is independent of the public ID — set it
with `UploadParams.AssetFolder` and search it with `asset_folder:`.

## Tags and contextual metadata

Tags can be set at upload or managed afterwards through the uploader:

```go
cld.Upload.AddTag(ctx, uploader.AddTagParams{Tag: "featured", PublicIDs: api.CldAPIArray{"id"}})
cld.Upload.RemoveTag(ctx, uploader.RemoveTagParams{Tag: "featured", PublicIDs: api.CldAPIArray{"id"}})
cld.Upload.ReplaceTag(ctx, uploader.ReplaceTagParams{Tag: "new", PublicIDs: api.CldAPIArray{"id"}})
cld.Upload.AddContext(ctx, uploader.AddContextParams{
	Context:   api.CldAPIMap{"alt": "Sample image"},
	PublicIDs: api.CldAPIArray{"id"},
})
```

For typed, validated fields use structured metadata instead — see
[Use structured metadata](use-structured-metadata.md).

## Rate limits

Admin and Search calls are rate-limited per hour. The design that keeps you well inside the
limit is to keep Admin calls out of request paths and cache what you read — treat them as
management operations rather than per-request lookups.

If you do exhaust the limit, `result.Error.Message` reports `Rate limit exceeded`; retry with
backoff. Note that an unsubscribed add-on can also surface as a rate-limit error rather than
a permission error.

## Troubleshooting

- Zero results from an expression you expected to match — see
  [Writing expressions that match](#writing-expressions-that-match). An unknown field is not
  an error; it simply matches nothing.
- `Query Error (at position 1)` — a leading wildcard or bare `*`.
- `Resource not found` for an asset you know exists — wrong `AssetType` (defaults to image),
  or you passed an asset ID where a public ID is required.
- An asset you just uploaded is missing from search results — index lag; look it up by ID.
- A delete reported success but removed nothing, or removed too much — API rejections and
  no-op matches both return `err == nil`. Check `result.Error.Message` and the returned
  `Deleted` map.

## Related

- [Use structured metadata](use-structured-metadata.md)
- [Handle errors](handle-errors.md)
- [Search expression reference](https://cloudinary.com/documentation/search_expressions.md)
- [Asset administration guide](https://cloudinary.com/documentation/go_asset_administration.md)
