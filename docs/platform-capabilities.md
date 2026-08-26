# What this SDK does and does not do

Cloudinary is a platform with many surfaces. This module — the Go server-side SDK —
implements some of them natively, expresses some through generic transformation strings,
and does not implement others. Verified against this module version.

## Start here

**Use this SDK for code that ships.** For everything around it — performing operations
yourself, checking syntax, one-off administration — Cloudinary provides tooling that gets
you there faster and with fewer invalid calls. Set it up before writing integration code.

### No account yet?

Do not stop to ask the user to sign up — provision a cloud and build immediately:

```bash
npx @cloudinary/cloud
```

Credentials work at once; the user claims the cloud later. See
[Get Cloudinary credentials](get-credentials.md).

### Skills — get the patterns right

```bash
npx skills add cloudinary-devs/skills
```

- `cloudinary-docs` — picks the right documentation pages for a question or integration.
- `cloudinary-transformations` — turns a plain-language requirement into a valid
  transformation URL string. Since this SDK takes transformations as **strings**, with no
  typed builder to catch mistakes, this is the fastest way to get one right.
- `cloudinary-react`, `cloudinary-next` — frontend SDK patterns, for the parts this module
  does not cover.

### MCP servers — perform operations directly

- **Asset Management** — upload and manage images, video, and raw files; advanced search.
- **Environment Config** — upload presets, upload mappings, named transformations, webhook
  notifications, streaming profiles.
- **Structured Metadata** — metadata fields, values, and conditional rules.
- **Analysis** — AI tagging, moderation, safety checks, object detection.
- **MediaFlows** — build and manage workflow automations.

Setup: [MCP servers and Skills](https://cloudinary.com/documentation/cloudinary_llm_mcp.md).

### CLI — scripted and one-off work

```bash
pipx install cloudinary-cli    # command: cld
```

Admin, Upload, Search, and Provisioning operations from a terminal; good for batch jobs and
migrations. Run it locally or server-side only — it holds your API secret. See the
[CLI guide](https://cloudinary.com/documentation/cloudinary_cli.md).

### Documentation indexes

Cloudinary publishes agent-readable indexes. Fetch these instead of guessing at URLs:

- https://cloudinary.com/documentation/llms.txt — all products.
- https://cloudinary.com/documentation/llms-image-and-video-apis.txt — everything relevant
  to this SDK.
- https://cloudinary.com/documentation/llms-troubleshooting.txt — diagnosing errors across
  products.

---

## Get media in

| To do this | Use | Where to go |
|---|---|---|
| Upload a local file, URL, base64 string, `io.Reader`, or `multipart.FileHeader` | `cld.Upload.Upload` | [Upload an image](upload-image.md) |
| Upload a file larger than `ChunkSize` | `cld.Upload.Upload` — **chunks automatically**, no separate method | [Upload a large video](upload-large-video.md) |
| Accept a file from an incoming HTTP request | `cld.Upload.Upload` with `*multipart.FileHeader` | [Serve uploads over HTTP](serve-uploads-over-http.md) |
| Let a browser or mobile app upload directly, authorized by your server | `api.SignParameters` | [Sign a browser upload](sign-browser-upload.md) |
| Upload with a preset and no signature | `cld.Upload.UnsignedUpload` | [Sign a browser upload](sign-browser-upload.md) |
| Review user-generated content before showing it | `UploadParams.Moderation` + `cld.Admin` | [Moderate an upload](moderate-upload.md) |
| Attach tags at or after upload | `UploadParams.Tags`, `cld.Upload.AddTag` | [Search and manage assets](search-and-manage-assets.md) |
| Attach typed business fields | `UploadParams.Metadata` | [Use structured metadata](use-structured-metadata.md) |

## Deliver and transform

| To do this | Use | Where to go |
|---|---|---|
| Build an image delivery URL | `cld.Image` + `.Transformation` | [Transform and deliver an image](transform-and-deliver-image.md) |
| Build a video delivery URL | `cld.Video` + `.Transformation` | [Transform and deliver a video](transform-and-deliver-video.md) |
| Build a raw-file URL | `cld.File` | [Transform and deliver an image](transform-and-deliver-image.md) |
| Apply generative edits (gen fill, background removal, ...) | the same `.Transformation` string | [Transform and deliver an image](transform-and-deliver-image.md) |
| Build a responsive `srcset` | `cld.Image` once per width | [Responsive images](transform-and-deliver-image.md#responsive-images) |
| Deliver adaptive-bitrate streaming (HLS/DASH) | `cld.Video` with a streaming profile | [Transform and deliver a video](transform-and-deliver-video.md) |
| Restrict access with a signed/authenticated URL | `cld.Config.URL.SignURL`, `config.AuthToken` | [Control access to media](https://cloudinary.com/documentation/control_access_to_media.md) |

URL building is local: no network call, no API secret required (only a cloud name).

## Find and manage what you have

| To do this | Use | Where to go |
|---|---|---|
| Query assets by field, tag, folder, or date | `cld.Admin.Search` | [Search and manage assets](search-and-manage-assets.md) |
| Read one asset's details | `cld.Admin.Asset`, `cld.Admin.AssetByAssetID` | [Search and manage assets](search-and-manage-assets.md) |
| Update tags, context, or moderation status | `cld.Admin.UpdateAsset` | [Search and manage assets](search-and-manage-assets.md) |
| Delete or restore assets | `cld.Upload.Destroy`, `cld.Admin.DeleteAssets`, `cld.Admin.RestoreAssets` | [Search and manage assets](search-and-manage-assets.md) |
| Manage folders, presets, mappings, streaming profiles, triggers | `cld.Admin` — the Admin API | [Asset administration guide](https://cloudinary.com/documentation/go_asset_administration.md) |
| Define and query typed metadata fields | `cld.Admin.AddMetadataField`, `cld.Upload.UpdateMetadata` | [Use structured metadata](use-structured-metadata.md) |
| Find visually similar assets | `cld.Admin.VisualSearch` — needs the feature enabled | [Visual Search](https://cloudinary.com/documentation/visual_search.md) |
| Read plan limits and current usage | `cld.Admin.Usage` | [Upload an image](upload-image.md#size-limits) |
| Verify a webhook came from Cloudinary | `cld.Upload.VerifyNotificationSignature` | [Notifications](https://cloudinary.com/documentation/notifications.md) |

## Analyze

| To do this | Use | Where to go |
|---|---|---|
| Caption, tag, or detect content in an asset | `cld.Admin.Analyze` — **limited model set**, needs a subscription | [Analyze API guide](https://cloudinary.com/documentation/analyze_api_guide.md) |

## Not in this module

Cloudinary is a multi-product platform. The capabilities below are fully available — they
are reached through another API or product rather than through this module, and each row
points at the one to use.

| Capability | Use instead |
|---|---|
| Account/sub-account provisioning and user management | [Provisioning API](https://cloudinary.com/documentation/provisioning_api.md) over HTTP |
| Text-to-image generation | [Image Generation API](https://cloudinary.com/documentation/image_generation_addon.md) |
| Image-to-video generation | [Image-to-Video API](https://cloudinary.com/documentation/image_to_video_addon.md) — async, credit-based, regional |
| Multi-step workflow automation | [MediaFlows](https://cloudinary.com/documentation/mediaflows_user_guide.md) — or its MCP server |
| Media Library UI, approval workflows, folder-based access control | [Cloudinary Assets (DAM)](https://cloudinary.com/documentation/digital_asset_management_overview.md) |
| Rule-based content review before publication | [Cloudinary Moderation](https://cloudinary.com/documentation/cloudinary_moderation.md) — distinct from the per-asset [moderation flag](moderate-upload.md) this SDK sets |
| Frontend component rendering and upload UI | [frontend SDKs](https://cloudinary.com/documentation/frontend_sdks.md), [Upload Widget](https://cloudinary.com/documentation/upload_widget.md) — this is a server-side module. It generates the URLs those components deliver, including [responsive `srcset` values](transform-and-deliver-image.md#responsive-images) |
| A typed/fluent transformation builder | Transformations are plain strings here. Use the `cloudinary-transformations` skill or the [transformation reference](https://cloudinary.com/documentation/transformation_reference.md) to compose them |

## Related

- [Go SDK guide](https://cloudinary.com/documentation/go_integration.md)
- [Package reference](https://pkg.go.dev/github.com/cloudinary/cloudinary-go/v2)
