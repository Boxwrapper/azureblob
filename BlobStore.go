package azureblob

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// IBlobStore interface for a BlobStore
type IBlobStore interface {
	Connect(container string) (azblob.ContainerURL, context.Context, error)
}

// BlobStore structural definition
type BlobStore struct {
	IBlobStore
	Name string
	Key  string
}

// Connect connects to an Azure Blob-Container of this Azure Storage-Account
func (store BlobStore) Connect(container string) (azblob.ContainerURL, context.Context, error) {
	credential, err := azblob.NewSharedKeyCredential(store.Name, store.Key)
	if err != nil {
		return azblob.ContainerURL{}, nil, fmt.Errorf("Invalid credentials with error: %s", err.Error())
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", store.Name, container))
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	containerContext := context.Background()
	_, err = containerURL.Create(containerContext, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil && err.(azblob.StorageError).ServiceCode() != azblob.ServiceCodeContainerAlreadyExists {
		return containerURL, containerContext, fmt.Errorf("Obtained container with error: %s", err.Error())
	}

	return containerURL, containerContext, nil
}
