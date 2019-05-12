# Azure Blobstorage

This three file library aims to make the usage of the [Azure SDK (Go) for Blobs](https://github.com/Azure/azure-storage-blob-go) a whole lot simpler and can be used as a starting point for much more complex projects. It was created for interfacing with Blobs on Azure Storageaccounts containing JSON data, though custom services can be created that are the able to store different types of documents on an Azure Storageaccount.

There are three components to this library:

- The `IBlobStore` abstracts the establishment of a connection to the Storageaccount.
- The `IBlobRepository` abstracts CRUD methods for interacting with the Blob data.
- The `IBlobService` for interfacing with the stored JSON, or mapping objects from to a Blob.

The blob service contained within the library is specialized to use JSON a storage format, but offerers a general purpose interface **IBlobService** from which other custom services can be derived, if necessary to your specific project.

## Example

This example demonstrates how the provided components can be used to retrieve and store data in Azure. It first uploads an `example` struct to the provided Storageaccount, then downloads and prints the retrived data to the standard output.

```go
package main

import (
    "azureblob"
    "fmt"
)

func main() {
    // Define the data that is stored
    type example struct {
        Message string `json:"message"`
    }

    // Authenticate with the Azure Storageaccount
    storageAccount := azureblob.BlobStore{
        Name: "<Storageaccount name>",
        Key:  "<Storageaccount token>"}

    // Retrive or create the blob container
    exampleBlob := azureblob.GetBlobService(storageAccount, "example")

    // Create a new entry
    inputData := example{"Hello World"}
    id, err := exampleBlob.Create(inputData)
    if err != nil {
        panic(err)
    }

    // Retrieve the stored data
    outputData := example{}
    err = exampleBlob.Read(id, &outputData)
    if err != nil {
        panic(err)
    }

    // Print the retrived data
    fmt.Println(outputData) // "{Hello World}"
}
```

## License

MIT License

Copyright &copy; 2019 - Samuel Oechsler

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
