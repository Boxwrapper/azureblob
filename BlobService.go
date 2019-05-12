package azureblob

import (
	"encoding/json"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
)

// IBlobService interface for a blobService
type IBlobService interface {
	Create(object interface{}) (string, error)
	List() ([]string, error)
	Read(id string, object interface{}) error
	Update(object interface{}) error
	Delete(id string) error
}

// blobService structural defintion
type blobService struct {
	IBlobService
	repo IBlobRepository
}

// GetBlobService retrives the blobService instance
func GetBlobService(blobStore IBlobStore, containerName string) IBlobService {
	instance := new(blobService)
	instance.repo = BlobRepository{
		BlobStore:     blobStore,
		ContainerName: containerName}

	return instance
}

// Create abstracts the process of creating a blob from an interface
func (service blobService) Create(object interface{}) (string, error) {
	generator, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	mappedObject := make(map[string]interface{})
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &mappedObject); err != nil {
		return "", err
	}
	identifier := generator.String()
	mappedObject["id"] = identifier

	data, err = json.Marshal(mappedObject)
	if err != nil {
		return identifier, err
	}

	if err := service.repo.Create(identifier, data); err != nil {
		return identifier, err
	}

	return identifier, nil
}

// List abstracts the process of listing blobs from a repository / provider
func (service blobService) List() ([]string, error) {
	entries, err := service.repo.List()
	if err != nil {
		return entries, err
	}
	return entries, nil
}

// Read abstracts the process of reading a blob and mapping it to an interface
func (service blobService) Read(id string, object interface{}) error {
	val := reflect.ValueOf(object)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("Argument 'object' must be passed by reference")
	}

	data, err := service.repo.Read(id)
	if err != nil && len(data) > 0 {
		return err
	} else if len(data) <= 0 {
		return fmt.Errorf("Blob with identifier '%s' does not exist", id)
	}

	mappedObject := make(map[string]interface{})
	if err := json.Unmarshal(data, &mappedObject); err != nil {
		return err
	}
	identifier := mappedObject["id"].(string)
	if identifier != id {
		return fmt.Errorf("Checksum '%s' does not match with the provided '%s'", identifier, id)
	}

	if err := json.Unmarshal(data, object); err != nil {
		return err
	}

	return nil
}

// Update abstracts the process of updating a blob from an interface
func (service blobService) Update(object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}

	mappedObject := make(map[string]interface{})
	if err := json.Unmarshal(data, &mappedObject); err != nil {
		return err
	}
	identifier := mappedObject["id"].(string)

	if exists, _ := service.repo.Read(identifier); len(exists) <= 0 {
		return fmt.Errorf("Blob with identifier '%s' does not exist", identifier)
	}

	if err := service.repo.Update(identifier, data); err != nil {
		return err
	}
	return nil
}

// Delete abstracts the process of deleting a blob on a repository / provider
func (service blobService) Delete(id string) error {
	if exists, _ := service.repo.Read(id); len(exists) <= 0 {
		return fmt.Errorf("Blob with identifier '%s' does not exist", id)
	}

	if err := service.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
