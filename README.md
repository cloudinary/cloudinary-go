Cloudinary Go SDK
==========

Cloudinary is a cloud service that offers a solution to a web application's entire image management pipeline.

Easily upload images to the cloud. Automatically perform smart image resizing, cropping and conversion without
installing any complex software. Integrate Facebook or Twitter profile image extraction in a snap, in any dimension and
style to match your website's graphics requirements. Images are seamlessly delivered through a fast CDN, and much much
more.

Cloudinary offers comprehensive APIs and administration capabilities and is easy to integrate with any web application,
existing or new.

Cloudinary provides URL and HTTP based APIs that can be easily integrated with any Web development framework.

For Go, Cloudinary provides a module for simplifying the integration even further.

## Setup ######################################################################

To install Cloudinary Go SDK, use `go get`:

```
go get github.com/cloudinary/cloudinary-go
```

## Usage

### Configuration

Each request for building a URL of a remote cloud resource must have the `CloudName` parameter set. Each request to our
secure APIs (e.g., image uploads, eager sprite generation) must have the `ApiKey` and `ApiSecret` parameters set.
See [API, URLs and access identifiers](http://cloudinary.com/documentation/api_and_access_identifiers) for more details.

Setting the `CloudName`, `ApiKey` and `ApiSecret` parameters can be done by initializing the Cloudinary object, or by
using the CLOUDINARY_URL environment variable / system property.

The entry point of the library is the Cloudinary struct.

    cld, _ := cloudinary.Create()

Here's an example of setting the configuration parameters programatically:

    cld, _ := cloudinary.CreateFromParams('n07t21i7','123456789012345', 'abcdeghijklmnopqrstuvwxyz12')

Another example of setting the configuration parameters by providing the CLOUDINARY_URL value:

    cld, _ := cloudinary.CreateFromUrl('cloudinary://123456789012345:abcdeghijklmnopqrstuvwxyz12@n07t21i7')

### Upload

Assuming you have your Cloudinary configuration parameters defined (`CloudName`, `ApiKey`, `ApiSecret`), uploading to
Cloudinary is very simple.

The following example uploads a local JPG to the cloud:

```
resp, err := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{});
```

The uploaded image is assigned a randomly generated public ID. The image is immediately available for a download through
a CDN:

```
println(resp.SecureURL)

// https://res.cloudinary.com/demo/image/upload/abcfrmo8zul1mafopawefg.jpg
```

You can also specify your own public ID:

```
resp, err := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{PublicID: "sample_remote.jpg"});
if err != nil {...}
println(resp.SecureURL)

// https://res.cloudinary.com/demo/image/upload/sample_remote.jpg
```

## Additional resources ##########################################################

Additional resources are available at:

* [Website](https://cloudinary.com)
* [Interactive demo](https://demo.cloudinary.com/default)
* [Knowledge Base](https://support.cloudinary.com/hc/en-us)
* [Documentation](https://cloudinary.com/documentation)
* [Upload API documentation](https://cloudinary.com/documentation/upload_images)

## Support

You can [open an issue through GitHub](https://github.com/cloudinary/cloudinary-go/issues).

Stay tuned for updates, tips and tutorials: [Blog](https://cloudinary.com/blog)
, [Twitter](https://twitter.com/cloudinary), [Facebook](https://www.facebook.com/Cloudinary).

## Join the Community ###########################################################

Impact the product, hear updates, test drive new features and more!
Join [here](https://www.facebook.com/groups/CloudinaryCommunity).

## Staying up to date  ##########################################################

To update Cloudinary Go SDK to the latest version, use `go get -u github.com/cloudinary/cloudinary-go`.

## Contributing  ################################################################

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.

## License ######################################################################

Released under the MIT license. 
