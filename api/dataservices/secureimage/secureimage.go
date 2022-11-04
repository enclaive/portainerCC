package secureimage

import (
	"fmt"

	"github.com/rs/zerolog/log"

	portainer "github.com/portainer/portainer/api"
)

const (
	BucketName = "secureimages"
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

func (service *Service) Create(secImgObject *portainer.SecureImage) error {

	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			secImgObject.ID = portainer.SecureImageID(id)
			return int(id), secImgObject
		},
	)
}

func (service *Service) SecureImages() ([]portainer.SecureImage, error) {
	var images = make([]portainer.SecureImage, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.SecureImage{},
		func(obj interface{}) (interface{}, error) {
			img, ok := obj.(*portainer.SecureImage)
			if !ok {
				log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("failed to convert to secure image object")
				return nil, fmt.Errorf("Failed to convert to secure image object: %s", obj)
			}

			images = append(images, *img)

			return &portainer.SecureImage{}, nil
		})

	return images, err
}
