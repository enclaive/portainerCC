package key

import (
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

func (service *Service) Create(keyObject *portainer.Key, data string) error {

	//only generate if not set
	// if keyObject.Key == nil {
	// 	// generate new rsa key
	// 	privatekey, err := GenerateMultiPrimeKeyForSGX(rand.Reader, 2, 3072)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Cannot generate RSA key\n")
	// 		return errors.New("Could not generate Key")
	// 	}

	// 	keyObject.Key = privatekey
	// }

	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			keyObject.ID = portainer.KeyID(id)
			return int(id), keyObject
		},
	)
}
