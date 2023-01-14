package confidentialtemplate

import (
	"fmt"

	"github.com/rs/zerolog/log"

	portainer "github.com/portainer/portainer/api"
)

const (
	BucketName = "confidentialtemplates"
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

func (service *Service) Create(conftemplateObject *portainer.ConfidentialTemplate) error {

	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			conftemplateObject.ID = portainer.ConfidentialTemplateId(id)
			return int(id), conftemplateObject
		},
	)
}

func (service *Service) ConfidentialTemplate(ID portainer.ConfidentialTemplateId) (*portainer.ConfidentialTemplate, error) {
	var template portainer.ConfidentialTemplate
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (service *Service) ConfidentialTemplates() ([]portainer.ConfidentialTemplate, error) {
	var templates = make([]portainer.ConfidentialTemplate, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.ConfidentialTemplate{},
		func(obj interface{}) (interface{}, error) {
			template, ok := obj.(*portainer.ConfidentialTemplate)
			if !ok {
				log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("failed to convert to confidential image object")
				return nil, fmt.Errorf("Failed to convert to confidential image object: %s", obj)
			}

			templates = append(templates, *template)

			return &portainer.ConfidentialTemplate{}, nil
		})

	return templates, err
}
