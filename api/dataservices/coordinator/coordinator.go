package coordinator

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
)

const (
	BucketName = "coordinator"
)

type Service struct {
	connection portainer.Connection
}

func (service *Service) BucketName() string {
	return BucketName
}

// NewService creates a new instance of this coordinator service.
func NewService(connection portainer.Connection) (*Service, error) {
	err := connection.SetServiceName(BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		connection: connection,
	}, nil
}

// Coordinators return an array containing all the coordinators
func (service *Service) Coordinators() ([]portainer.Coordinator, error) {
	var coordinators = make([]portainer.Coordinator, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.Coordinator{},
		func(obj interface{}) (interface{}, error) {
			coordinator, ok := obj.(*portainer.Coordinator)
			if !ok {
				logrus.WithField("obj", obj).Errorf("Failed to convert to Coordinator object")
				return nil, fmt.Errorf("Failed to convert to Coordinator object: %s", obj)
			}

			coordinators = append(coordinators, *coordinator)
			return &portainer.Coordinator{}, nil
		})
	return coordinators, err
}

// Coordinator returns the coordinator with the specified id
func (service *Service) Coordinator(ID portainer.CoordinatorID) (*portainer.Coordinator, error) {

	var coordinator portainer.Coordinator
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &coordinator)
	if err != nil {
		return nil, err
	}

	return &coordinator, nil
}

// Create creates a new coordinator
func (service *Service) Create(coordinatorObject *portainer.Coordinator) error {
	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			coordinatorObject.ID = portainer.CoordinatorID(id)
			return int(id), coordinatorObject
		},
	)
}

// Update an existing coordinator
func (service *Service) Update(ID portainer.CoordinatorID, keyObject *portainer.Coordinator) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.UpdateObject(BucketName, identifier, keyObject)
}

// Remove an existing coordinator
func (service *Service) Delete(ID portainer.CoordinatorID) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.DeleteObject(BucketName, identifier)
}
