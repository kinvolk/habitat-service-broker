package broker

import (
	"fmt"
	"sync"

	habv1beta1 "github.com/habitat-sh/habitat-operator/pkg/apis/habitat/v1beta1"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/osb-broker-lib/pkg/broker"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewBrokerLogic is a hook that is called with the Options the program is run with.
func NewBrokerLogic(o *Options, client *Client) (*BrokerLogic, error) {
	return &BrokerLogic{
		async:      o.Async,
		KubeClient: client,
	}, nil
}

// BrokerLogic provides an implementation of the broker.BrokerLogic interface.
type BrokerLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	KubeClient *Client
}

// Client stores all the information specfic to Kubernetes.
type Client struct {
	KubeClient kubernetes.Interface
	Client     *rest.RESTClient
}

var _ broker.Interface = &BrokerLogic{}

func (b *BrokerLogic) GetCatalog(c *broker.RequestContext) (*osb.CatalogResponse, error) {
	response := &osb.CatalogResponse{}

	// TODO (lilic): At some point move these at a more appropriate place.
	data := `
---
services:
- name: nginx-habitat
  id: 1ac7de1d-d89a-41c7-b9a8-744f9256e375
  description: Nginx packaged with Habitat
  bindable: false
  plan_updateable: false
  metadata:
    displayName: "Habitat Nginx service"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: 86064792-7ea2-467b-af93-ac9694d96d5b
    description: The default plan for the Nginx Habitat service
    free: true
    schemas:
      service_instance:
        create:
          "$schema": "http://json-schema.org/draft-04/schema"
          "type": "object"
          "title": "Parameters"
          "properties":
          - "name":
              "title": "Some Name"
              "type": "string"
              "maxLength": 63
              "default": "My Name"
          - "color":
              "title": "Color"
              "type": "string"
              "default": "Clear"
              "enum":
              - "Clear"
              - "Beige"
              - "Grey"
- name: redis-habitat
  id: 50e86479-4c66-4236-88fb-a1e61b4c9448 
  description: Redis packaged with Habitat
  bindable: false
  plan_updateable: false
  metadata:
    displayName: "Habitat Redis service"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: 002341cf-f895-49f4-ba04-bb70291b895c
    description: The default plan for the Redis Habitat example service
    free: true
    schemas:
      service_instance:
        create:
          "$schema": "http://json-schema.org/draft-04/schema"
          "type": "object"
          "title": "Parameters"
          "properties":
          - "name":
              "title": "Some Name"
              "type": "string"
              "maxLength": 63
              "default": "My Name"
          - "color":
              "title": "Color"
              "type": "string"
              "default": "Clear"
              "enum":
              - "Clear"
              - "Beige"
              - "Grey"
`

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BrokerLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*osb.ProvisionResponse, error) {
	b.Lock()
	defer b.Unlock()

	response := osb.ProvisionResponse{}

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	fmt.Println("Started provisioning...")

	hab, err := generateHabitatObject(request.PlanID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = b.createHabitatResource(hab)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// TODO: Figure out what else needs to be created here for each of the pre-packaged services.
	// TODO: when statefulsets are merged into hab-operator, create a PV/PVC here as well.

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*osb.DeprovisionResponse, error) {
	b.Lock()
	defer b.Unlock()

	response := osb.DeprovisionResponse{}

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*osb.LastOperationResponse, error) {
	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*osb.BindResponse, error) {
	b.Lock()
	defer b.Unlock()

	response := osb.BindResponse{}

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &osb.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*osb.UpdateInstanceResponse, error) {
	response := osb.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}
