package key

import (
	"fmt"

	"github.com/rs/zerolog/log"

	portainer "github.com/portainer/portainer/api"
)

const (
	BucketName = "keys"
)

type Service struct {
	connection portainer.Connection
}

func (service *Service) BucketName() string {
	return BucketName
}

// NewService creates a new instance of this conf. compute service.
func NewService(connection portainer.Connection) (*Service, error) {
	err := connection.SetServiceName(BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		connection: connection,
	}, nil
}

func (service *Service) Create(keyObject *portainer.Key) error {

	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			keyObject.ID = portainer.KeyID(id)
			return int(id), keyObject
		},
	)
}

func (service *Service) Key(ID portainer.KeyID) (*portainer.Key, error) {
	var key portainer.Key
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (service *Service) Keys() ([]portainer.Key, error) {
	var keys = make([]portainer.Key, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.Key{},
		func(obj interface{}) (interface{}, error) {
			key, ok := obj.(*portainer.Key)
			if !ok {
				log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("failed to convert to Key object")
				return nil, fmt.Errorf("Failed to convert to Key object: %s", obj)
			}

			keys = append(keys, *key)

			return &portainer.Key{}, nil
		})

	return keys, err
}

func (service *Service) Update(ID portainer.KeyID, keyObject *portainer.Key) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.UpdateObject(BucketName, identifier, keyObject)
}

func (service *Service) Delete(ID portainer.KeyID) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.DeleteObject(BucketName, identifier)
}
