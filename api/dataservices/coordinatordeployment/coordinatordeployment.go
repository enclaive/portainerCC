package coordoinatordeployment

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
)

const (
	BucketName = "coordinatordeployment"
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
func (service *Service) CoordinatorDeployments() ([]portainer.CoordinatorDeployment, error) {
	var coordinatorDeployments = make([]portainer.CoordinatorDeployment, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.CoordinatorDeployment{},
		func(obj interface{}) (interface{}, error) {
			coordinatorDeployment, ok := obj.(*portainer.CoordinatorDeployment)
			if !ok {
				logrus.WithField("obj", obj).Errorf("Failed to convert to CoordinatorDeployment object")
				return nil, fmt.Errorf("Failed to convert to CoordinatorDeployment object: %s", obj)
			}

			coordinatorDeployments = append(coordinatorDeployments, *coordinatorDeployment)
			return &portainer.Coordinator{}, nil
		})
	return coordinatorDeployments, err
}

// Coordinator returns the coordinator with the specified id
func (service *Service) CoordinatorDeployment(ID portainer.CoordinatorDeploymentID) (*portainer.CoordinatorDeployment, error) {

	var coordinatorDeployment portainer.CoordinatorDeployment
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &coordinatorDeployment)
	if err != nil {
		return nil, err
	}

	return &coordinatorDeployment, nil
}

// Create creates a new coordinator
func (service *Service) Create(coordinatorDeploymentObject *portainer.CoordinatorDeployment) error {
	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			coordinatorDeploymentObject.ID = portainer.CoordinatorDeploymentID(id)
			return int(id), coordinatorDeploymentObject
		},
	)
}

// Update an existing coordinator deployment
func (service *Service) Update(ID portainer.CoordinatorDeploymentID, keyObject *portainer.CoordinatorDeployment) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.UpdateObject(BucketName, identifier, keyObject)
}

// Remove an existing coordinator
func (service *Service) Delete(ID portainer.CoordinatorDeploymentID) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.DeleteObject(BucketName, identifier)
}
