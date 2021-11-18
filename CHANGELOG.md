1.4.0 / 2021-11-18
==================

New Functionality And Features
------------------------------
  * Add support for `Moderation` in `UploadResult`

Other Changes
-------------
  * Fix `UpdateMetadata` Upload API parameters

1.3.0 / 2021-08-20
==================

New Functionality And Features
------------------------------
  * Add support for query parameters in `CLOUDINARY_URL`
  * Add support for `DeliveryType` in `Assets` Admin API
  * Add support for `ReorderMetadataFieldDatasource` Admin API
  * Add `Info` field to `AssetResult`

Other Changes
-------------
* Add code generation for typed setters
* Add support for acceptance tests

1.2.0 / 2021-06-24
==================

New Functionality And Features
------------------------------
  * Add support for URL generation
  * Add missing `NextCursor` to `DeleteAssetsResult`

Other Changes
-------------
  * Fix request signature of the Upload API

1.1.0 / 2021-05-13
==================

New Functionality And Features
------------------------------
  * Add `Logger`
  * Add `NextCursor` field to search API result

Other Changes
-------------
  * Fix Upload with Context

1.0.0 / 2021-04-06
==================

Breaking Changes
------------------------------
  * Rename `Create()` to `New()`
  * Rename `NewFromUrl` to `NewFromURL`
  * Rename `Account` config to `Cloud`
  * Apply `golint` fixes
  * Unexport redundant public constants
  * Refactor Search Api

New Functionality And Features
------------------------------
  * Add `PrivateDownloadURL` helper
  * Add support for `Eager` upload transformations
  * Add support for large file upload
  * Add SDK usage example
  * Add `CONTRIBUTING.md`
    
Other Changes
-------------
  * Improve error handling in Configuration
  * Fix `Config.Api.Timeout` type
  * Fix api call timeout handling
  * Fix formatting
  * Update README
  * Add travis support
  * Refactor tests
  * Code cleanup
  * Add a link to the package documentation
  * Update README code samples


0.2.0 / 2021-01-05
==================

  * Fix package name
  * Reorganise files
  * Cleanup files

0.1.0 / 2021-01-04
========================

  * The first public beta
