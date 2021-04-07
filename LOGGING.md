### Logging in Cloudinary SDK
The default logger in Cloudinary Go SDK is `go log`.

You can use any log library by overwriting standard SDK logging functions.

#### Using logrus with the SDK
```go
package main

import (
	"github.com/cloudinary/cloudinary-go"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.New()
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	// Initialize your logger somewhere in your code.
	// Set cloudinary.Logger.ErrorLogger and cloudinary.Logger.DebugLogger fields with logrus functions
	var logger = logrus.New()
	cld.Logger.ErrorLogger = logger.WithField("source", "cloudinary").Errorln
	cld.Logger.DebugLogger = logger.WithField("source", "cloudinary").Debugln
}
```

#### Using Zap with the SDK
```go
package main

import (
	"github.com/cloudinary/cloudinary-go"
	"go.uber.org/zap"
	"log"
)

func main() {
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.New()
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	// Initialize your logger somewhere in your code.
	// Set cloudinary.Logger.ErrorLogger and cloudinary.Logger.DebugLogger fields with zap sugared logger functions
	var zapLogger, _ = zap.NewDevelopment()
	var zapSugared = zapLogger.Sugar()

	cld.Logger.ErrorLogger = zapSugared.With("source", "cloudinary").Error
	cld.Logger.DebugLogger = zapSugared.With("source", "cloudinary").Debug
}
```

#### Logging level

You can change logging level with `Logger.SetLevel()` function.

Possible options:
- `logger.NONE`  - disabling logging from the SDK
- `logger.ERROR` - enable logging only for error messages
- `logger.DEBUG` - enable debug logs
