# sassy
> It's SASsy az ...

Sassy provides a Go library to generate Azure SAS tokens.

## Examples
### Storage
#### Generating an Account SAS

```go
package main

import (
	"fmt"
	"log"

	"github.com/matthewhartstonge/sassy/storage"
	"github.com/matthewhartstonge/sassy/storage/versions"
)

func main() {
	// Get SAS~sy with it...
	sas, err := storage.NewAccountSAS(
		"yourStorageAccountName",
		"yourStorageAccountKey",
		// signedVersion specifies what API version to use in order to generate
		// the storage SAS token - must be set to version 2015-04-05 or later.
		// Current valid versions can be found in `storage/versions/versions.go`
		versions.Latest.String(),
		// signedServices supports:
		// Blob  = "b"
		// Queue = "q"
		// Table = "t"
		// File  = "f"
		"bqtf",
		// signedResourceTypes supports:
		// Service   = "s"
		// Container = "c"
		// Object    = "o"
		"sco",
		// signedPermissions doesn't care what order you specify them in, it 
		// works out the required order for you, and negates the permissions 
		// not eligible based on your specified API version :3
		//
		// Note: 
		// - There are too many permissions to list here, so instead refer to:
		//   https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#permissions-for-a-directory-container-or-blob
		"rlw",
		// signedExpiry aims to be developer and ops friendly by first and 
		// foremost attempting to parse a given date in your local timezone, so
		// you no longer have to worry about converting to UTC - unless you 
		// want to... Then add a Z!
		// Formats supported:
		// - YYYY-MM-DD
		// - YYYY-MM-DD<TZDSuffix>
		// - YYYY-MM-DDThh:mm
		// - YYYY-MM-DDThh:mm<TZDSuffix>
		// - YYYY-MM-DDThh:mm:ss
		// - YYYY-MM-DDThh:mm:ss<TZDSuffix>
		"2021-12-12",
	)
	if err != nil {
		// You broke my SAS... :(
		log.Fatal(err)
	}

	// Get a signed SAS token:
	fmt.Println(sas.Token())
}
```

## TODO
* Storage: Service SAS generation 
* Storage: User Delegation SAS generation
* CLI tool
