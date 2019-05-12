package azureblob

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// IBlobRepository interface for a BlobRepository
type IBlobRepository interface {
	Create(fileName string, data []byte) error
	List() ([]string, error)
	Read(fileName string) ([]byte, error)
	Update(fileName string, data []byte) error
	Delete(fileName string) error
}

// BlobRepository structural definition
type BlobRepository struct {
	IBlobRepository
	BlobStore     IBlobStore
	ContainerName string
}

// Create crates a blob and uploads the provided data
func (repo BlobRepository) Create(fileName string, data []byte) error {
	return repo.Update(fileName, data)
}

// List lists the blobs contained in this BlobRepository
func (repo BlobRepository) List() ([]string, error) {
	containerURL, containerContext, err := repo.BlobStore.Connect(
		strings.ToLower(repo.ContainerName))
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(containerContext, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return fileNames, fmt.Errorf("Listing failed with error: %s", err.Error())
		}

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			fileNames = append(fileNames, blobInfo.Name)
		}
	}

	return fileNames, nil
}

// Read opens a blob and downloads its data
func (repo BlobRepository) Read(fileName string) ([]byte, error) {
	containerURL, containerContext, err := repo.BlobStore.Connect(
		strings.ToLower(repo.ContainerName))
	if err != nil {
		return nil, err
	}

	blobURL := containerURL.NewBlockBlobURL(fileName)
	var buffer []byte
	if err = downloadToBuffer(containerContext, blobURL, &buffer); err != nil {
		return buffer, err
	}
	return buffer, nil
}

// Update updates an existing blob and uploads the provided data
func (repo BlobRepository) Update(fileName string, data []byte) error {
	containerURL, containerContext, err := repo.BlobStore.Connect(
		strings.ToLower(repo.ContainerName))
	if err != nil {
		return err
	}

	blobURL := containerURL.NewBlockBlobURL(fileName)
	if err := uploadFromBuffer(containerContext, blobURL, data); err != nil {
		return fmt.Errorf("Upload failed with error: %s", err.Error())
	}
	return nil
}

// Delete deletes a specified blob and its associated data
func (repo BlobRepository) Delete(fileName string) error {
	containerURL, containerContext, err := repo.BlobStore.Connect(
		strings.ToLower(repo.ContainerName))
	if err != nil {
		return err
	}

	blobURL := containerURL.NewBlockBlobURL(fileName)
	if _, err = blobURL.Delete(containerContext, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{}); err != nil {
		return fmt.Errorf("Delete failed with error: %s", err.Error())
	}
	return nil
}

func uploadFromBuffer(containerContext context.Context, blobURL azblob.BlockBlobURL, buffer []byte) error {
	_, err := azblob.UploadBufferToBlockBlob(containerContext, buffer, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		return err
	}
	return nil
}

func downloadToBuffer(containerContext context.Context, blobURL azblob.BlockBlobURL, buffer *[]byte) error {
	response, err := blobURL.Download(containerContext, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return err
	}

	body := response.Body(azblob.RetryReaderOptions{
		MaxRetryRequests: 20})

	data := bytes.Buffer{}
	if _, err = data.ReadFrom(body); err != nil {
		return err
	}
	*buffer = data.Bytes()

	return nil
}
