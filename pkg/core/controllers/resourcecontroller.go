//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package controllers

import (
	"TFRP/pkg/core/apierror"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/entities"
	"TFRP/pkg/core/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

// GetResourceController returns a resource
func GetResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	resourcePackage := entities.ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(fullyQualifiedResourceID, &resourcePackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			err.Error())
		return
	}

	provider := engines.GetProvider(resourcePackage.ProviderType)

	cfg, err := config.Load(resourcePackage.Config)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			err.Error())
		return
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		err = provider.Configure(terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				err.Error())
			return
		}
	}

	info := &terraform.InstanceInfo{
		Type: resourcePackage.ResourceType,
	}

	// Call refresh
	resourceState, err := provider.Refresh(info, resourcePackage.State)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			err.Error())
		return
	}

	resourcePackage.State = resourceState

	// insert Document in collection
	err = storage.GetResourceDataProvider().Insert(&resourcePackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to insert data: %s", err))
		return
	}

	responseContent, err := json.Marshal(resourcePackage.ToDefinition())
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize response content: %s", err))
		return
	}

	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// PutResourceController creates/updates a resource
func PutResourceController(request *restful.Request, response *restful.Response) {
	resourceDefinition := entities.ResourceDefinition{}

	rawBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to read request content: %s", err))
		return
	}

	err = json.Unmarshal(rawBody, &resourceDefinition)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to deserialize request content: %s", err))
		return
	}

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err = storage.GetProviderRegistrationDataProvider().Find(resourceDefinition.Properties.ProviderID, &providerRegistrationPackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("The provider registration %s was not found.", resourceDefinition.Properties.ProviderID))
		return
	}

	resourceSpec, err := json.Marshal(resourceDefinition.Properties.Settings)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize resource property settings: %s", err))
		return
	}

	configFile := getConfigFileInJSON(
		providerRegistrationPackage.ProviderType,
		providerRegistrationPackage.Credentials,
		resourceDefinition,
		engines.GetResourceName(request), resourceSpec)
	fmt.Printf("%s", configFile)

	provider := engines.GetProvider(providerRegistrationPackage.ProviderType)

	cfg, err := config.Load(configFile)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Failed to parse config file: %s", err))
		return
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		err = provider.Configure(terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				fmt.Sprintf("Failed to init provider: %s", err))
			return
		}
	}

	info := &terraform.InstanceInfo{
		Type: resourceDefinition.Properties.ResourceType,
	}

	for _, v := range cfg.Resources {
		state := new(terraform.InstanceState)
		state.Init()
		diff, err := provider.Diff(info, state, terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				fmt.Sprintf("Failed to call provider diff: %s", err))
			return
		}

		// Call apply to create resource
		resourceState, err := provider.Apply(info, state, diff)
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				fmt.Sprintf("Failed to create resourse: %s", err))
			return
		}
		fmt.Printf("%s", resourceState.ID)

		fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

		// insert Document in collection
		err = storage.GetResourceDataProvider().Insert(&entities.ResourcePackage{
			ResourceID:   fullyQualifiedResourceID,
			StateID:      resourceState.ID,
			State:        resourceState,
			Config:       configFile,
			ResourceType: resourceDefinition.Properties.ResourceType,
			ProviderType: providerRegistrationPackage.ProviderType,
		})
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusInternalServerError,
				apierror.InternalError,
				apierror.InternalOperationError,
				fmt.Sprintf("Failed to insert data: %s", err))
			return
		}
	}

	responseContent, err := json.Marshal(resourceDefinition)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize response content: %s", err))
		return
	}

	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// DeleteResourceController deletes a resource
func DeleteResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	resourcePackage := entities.ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(fullyQualifiedResourceID, &resourcePackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			fmt.Sprintf("Resource with id '%s' was not found", fullyQualifiedResourceID))
		return
	}

	provider := engines.GetProvider(resourcePackage.ProviderType)

	cfg, err := config.Load(resourcePackage.Config)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			err.Error())
		return
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		err = provider.Configure(terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				fmt.Sprintf("Failed to init provider: %s", err))
			return
		}
	}

	info := &terraform.InstanceInfo{
		Type: resourcePackage.ResourceType,
	}

	diff := new(terraform.InstanceDiff)
	diff.Destroy = true

	// Call apply to delete resource
	resourceState, err := provider.Apply(info, resourcePackage.State, diff)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Failed to delete resourse: %s", err))
		return
	}

	if resourceState == nil {
		err := storage.GetResourceDataProvider().Remove(fullyQualifiedResourceID)
		if err != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusInternalServerError,
				apierror.InternalError,
				apierror.InternalOperationError,
				fmt.Sprintf("Failed to delete resource '%s' from storage: %s", fullyQualifiedResourceID, err))
			return
		}
	}

	response.WriteHeader(http.StatusOK)
}

func getConfigFileInJSON(providerType string, providerSpec []byte, resource entities.ResourceDefinition, resourceName string, resourceSpec []byte) string {
	return fmt.Sprintf(`
		{
			"provider": {
				"%s": %s
			},
			"resource": {
				"%s": {
					"%s": %s
				}
			}
		}
`, providerType, string(providerSpec), resource.Properties.ResourceType, resourceName, string(resourceSpec))
}
