# Moderate an upload

## When to use

Content uploaded by users must be reviewed before it is shown. Moderation in Cloudinary is
stateful: an asset carries a moderation status, and **your application** is responsible for
showing approved assets only.

## `pending` does not block delivery

**By default a moderated asset is deliverable from the moment it is uploaded.** Verified
against a live cloud: an asset with `moderation: manual` in `pending` status returns
**HTTP 200** on both its original URL and a transformed URL. Nothing 404s, nothing is
withheld.

The status is **metadata your application gates on**, which keeps the decision in your hands:
you choose what "approved enough to show" means for your product. Build the check into the
code path that renders user-generated content.

Blocking delivery of non-approved assets can also be configured at the product-environment
level — it is not an upload parameter, so contact Cloudinary support to enable it. Gate on
the status in your own code either way.

## Where the status lives

`UploadResult.Moderation` is a **slice**, since an asset can pass through several moderation
kinds in a chain:

```go
if len(result.Moderation) > 0 {
	fmt.Println(result.Moderation[0].Kind)   // "manual"
	fmt.Println(result.Moderation[0].Status) // api.Pending
}
```

Check the length first, as above: an asset uploaded without moderation has an empty slice.

`Status` is typed as `api.ModerationStatus`, with constants for the three statuses you act
on: `api.Pending`, `api.Approved`, and `api.Rejected`. The platform also reports **`queued`**
(waiting for an add-on to run) and **`aborted`** (an earlier moderation in a chain rejected
the asset) — compare those against the string:

```go
switch result.Moderation[0].Status {
case api.Pending, api.Approved, api.Rejected:
	// covered by constants
case "queued", "aborted":
	// reported by the platform; compare as strings
}
```

Note which type each call site takes — the compiler will tell you:

| Where | Type |
|---|---|
| `Moderation.Status` (result) | `api.ModerationStatus` |
| `admin.UpdateAssetParams.ModerationStatus` | `api.ModerationStatus` |
| `admin.AssetsByModerationParams.Status` | plain `string` — pass `string(api.Pending)` |

## Complete flow (manual review queue)

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const publicID = "examples/moderated-upload"

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

	// 1. Upload into the moderation queue — the asset starts as "pending".
	uploaded, err := cld.Upload.Upload(ctx, "https://res.cloudinary.com/demo/image/upload/sample.jpg",
		uploader.UploadParams{
			PublicID:   publicID,
			Moderation: "manual",
			Overwrite:  api.Bool(true),
		})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}
	if len(uploaded.Moderation) > 0 {
		fmt.Println("status:", uploaded.Moderation[0].Status) // pending
	}
	// NOTE: uploaded.SecureURL already serves this image. Do not treat it as gated.

	// 2. Your review UI lists the queue.
	// Note: Status here is a plain string, not api.ModerationStatus — the constants
	// are typed, so pass string(api.Pending) rather than api.Pending.
	queue, err := cld.Admin.AssetsByModeration(ctx, admin.AssetsByModerationParams{
		Kind:       "manual",
		Status:     string(api.Pending),
		MaxResults: 50,
	})
	if err != nil {
		return fmt.Errorf("list queue: %w", err)
	}
	if queue.Error.Message != "" {
		return fmt.Errorf("list queue rejected: %s", queue.Error.Message)
	}
	fmt.Printf("pending review: %d asset(s)\n", len(queue.Assets))

	// 3. A reviewer records the decision.
	updated, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID:         publicID,
		ModerationStatus: api.Approved,
	})
	if err != nil {
		return fmt.Errorf("update transport: %w", err)
	}
	if updated.Error.Message != "" {
		return fmt.Errorf("update rejected: %s", updated.Error.Message)
	}

	// 4. Only now does your application link to the asset.
	fmt.Println("approved:", updated.SecureURL)
	return nil
}
```

## Automatic moderation

Pass an add-on name instead of `manual` to get an automated verdict.

**Prerequisite — a human has to do this, not your code.** Every value below except `manual`
requires its add-on to be registered on the account first, from the
[Add-ons page](https://cloudinary.com/documentation/cloudinary_add_ons.md) in the console.
Some third-party add-ons also require reviewing and accepting the provider's terms of
service as part of registration. Neither step has an API. `manual` needs no add-on and no
terms accepted, which is why the flow above uses it.

| Value | Moderates | Add-on |
|---|---|---|
| `aws_rek` | images | Amazon Rekognition AI Moderation |
| `aws_rek_video` | video | Amazon Rekognition Video Moderation |
| `google_video_moderation` | video | Google AI Video Moderation |
| `webpurify` | images | WebPurify Image Moderation |
| `perception_point` | any asset | Perception Point Malware Detection |
| `duplicate:<threshold>` | images | Cloudinary Duplicate Image Detection |

Combine several with a pipe — the order is the order they run in, and `manual` must be last
(`"aws_rek|duplicate:0.9|manual"`). The first moderation starts as `pending` and the rest as
`queued`; if one rejects, the remaining become `aborted` and the asset's final status is
`rejected`. Always set a `NotificationURL` when requesting several, and verify the webhook
with `cld.Upload.VerifyNotificationSignature`.

An automated verdict can still be overridden by a human with `UpdateAsset` +
`ModerationStatus`.

## Design rules

- **Model moderation as a state machine, not a boolean**, and keep the pending state visible
  in your product (placeholder, "under review" label). Remember the URL works regardless, so
  the gate has to be in your code and in your data model.
- **Store `AssetID`, gate on your own copy of the status.** Re-reading Cloudinary on every
  page render is an Admin API call in a request path, which is rate-limited.
- **Keep human override even with automatic moderation** — machine verdicts are drafts for
  anything with legal or brand consequences.
- **Rejected assets stay in storage** unless you delete them; decide your retention policy.
- **Serve a placeholder for non-approved assets** rather than relying on the URL failing,
  because it will not fail.

## Troubleshooting

- A pending asset is publicly visible — expected. Nothing blocks delivery by default;
  enforcement is your application's responsibility. It is not an upload parameter: contact
  Cloudinary support to have it configured for your product environment.
- `Moderation <value> moderation is not valid` — the moderation value is misspelled; use one
  from the table above. Verified: this arrives with `err == nil`, in
  `result.Error.Message`.
- `You don't have an active subscription for <add-on>` — register the add-on in the console,
  and accept the provider's terms if it is third-party. Note that an unsubscribed add-on can
  also surface as a **rate-limit** error rather than a permission error, which is deeply
  unobvious.
- `index out of range` panic reading the status — the asset was uploaded without moderation,
  so `result.Moderation` is empty. Check `len` first.
- The status did not change after `UpdateAsset` — check `result.Error.Message`; an API
  rejection returns `err == nil`. See [Handle errors](handle-errors.md).

## Related

- Runnable example: `examples/moderate-upload/main.go` (in the repository)
- [Moderate assets](https://cloudinary.com/documentation/moderate_assets.md) — statuses,
  delivery behaviour, and the available moderation add-ons
- [Moderation guide](https://cloudinary.com/documentation/cloudinary_moderation.md) — the
  separate rule-based product, distinct from this per-asset flag
- [Handle errors](handle-errors.md)
