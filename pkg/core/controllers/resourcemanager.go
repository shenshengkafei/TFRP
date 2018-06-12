//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package controllers

import (
	"TFRP/pkg/core/apierror"
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/entities"
	"TFRP/pkg/core/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

// ResourceManager is the resource manager
type ResourceManager struct {
	BaseHandler
}

// NewResourceManager create a new resource manager
func NewResourceManager(
	providerRegistrationDataProvider *storage.ProviderRegistrationDataProvider,
	resourceDataProvider *storage.ResourceDataProvider) (resourceManager *ResourceManager) {
	resourceManager = new(ResourceManager)
	resourceManager.ProviderRegistrationDataProvider = providerRegistrationDataProvider
	resourceManager.ResourceDataProvider = resourceDataProvider
	return resourceManager
}

// GetResourceController returns a resource
func (resourceManager *ResourceManager) GetResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	resourcePackage := entities.ResourcePackage{}
	err := resourceManager.ResourceDataProvider.FindPackage(fullyQualifiedResourceID, &resourcePackage)
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

	if resourcePackage.State != nil {
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

		if resourceState == nil {
			err = resourceManager.ResourceDataProvider.RemovePackage(fullyQualifiedResourceID)
			apierror.WriteErrorToResponse(
				response,
				http.StatusNotFound,
				apierror.ClientError,
				apierror.NotFound,
				fmt.Sprintf("Resource with id '%s' was not found", fullyQualifiedResourceID))
			return
		}

		resourcePackage.State = resourceState

		// insert Document in collection
		err = resourceManager.ResourceDataProvider.InsertPackage(&resourcePackage)
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
func (resourceManager *ResourceManager) PutResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)
	resourceDefinition := entities.ResourceDefinition{}

	rawBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Request content is invalid: %s", err))
		return
	}

	err = json.Unmarshal(rawBody, &resourceDefinition)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Request content cannot be deserialized as JSON: %s", err))
		return
	}

	validationError := engines.ValidateResourceDefinition(&resourceDefinition)
	if validationError != nil {
		apierror.WriteErrorToResponseWitAPIError(
			response,
			http.StatusBadRequest,
			validationError)
		return
	}

	// Try to get provider registartion document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err = resourceManager.ProviderRegistrationDataProvider.FindPackage(resourceDefinition.Properties.ProviderID, &providerRegistrationPackage)
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
		providerRegistrationPackage.Settings,
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

	err = cfg.Validate()
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Invalid config file: %s", err))
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
		_, errs := provider.ValidateResource(resourceDefinition.Properties.ResourceType, terraform.NewResourceConfig(v.RawConfig))
		if errs != nil {
			apierror.WriteErrorToResponse(
				response,
				http.StatusBadRequest,
				apierror.ClientError,
				apierror.BadRequest,
				fmt.Sprintf("The resource settings are invalid: %s", errs))
			return
		}

		state := new(terraform.InstanceState)
		state.Init()

		// Get Document from collection
		resourcePackage := entities.ResourcePackage{}
		err := resourceManager.ResourceDataProvider.FindPackage(fullyQualifiedResourceID, &resourcePackage)
		if err == nil {
			if strings.EqualFold(resourcePackage.ProvisioningState, consts.ProvisioningStateAccepted) {
				apierror.WriteErrorToResponse(
					response,
					http.StatusConflict,
					apierror.ClientError,
					apierror.Conflict,
					fmt.Sprintf("Cannot create Resource with id '%s' as it is being provisioned", fullyQualifiedResourceID))
				return
			} else if resourcePackage.State != nil {
				// Call refresh
				state, err = provider.Refresh(info, resourcePackage.State)
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
		}

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

		// If we have no diff, we have nothing to do!
		if diff.Empty() {
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
			return
		}

		// insert Document in collection
		err = resourceManager.ResourceDataProvider.InsertPackage(&entities.ResourcePackage{
			ResourceID:        fullyQualifiedResourceID,
			ProvisioningState: consts.ProvisioningStateAccepted,
			Config:            configFile,
			ResourceType:      resourceDefinition.Properties.ResourceType,
			ProviderType:      providerRegistrationPackage.ProviderType,
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

		go func() {
			// Call apply to create resource
			resourceState, err := provider.Apply(info, state, diff)
			if err != nil {
				resourceManager.ResourceDataProvider.InsertPackage(&entities.ResourcePackage{
					ResourceID:               fullyQualifiedResourceID,
					ProvisioningState:        consts.ProvisioningStateFailed,
					ProvisioningErrorCode:    string(apierror.BadRequest),
					ProvisioningErrorMessage: err.Error(),
					Config:       configFile,
					ResourceType: resourceDefinition.Properties.ResourceType,
					ProviderType: providerRegistrationPackage.ProviderType,
				})
				return
			}

			// insert Document in collection
			err = resourceManager.ResourceDataProvider.InsertPackage(&entities.ResourcePackage{
				ResourceID:        fullyQualifiedResourceID,
				StateID:           resourceState.ID,
				State:             resourceState,
				ProvisioningState: consts.ProvisioningStateSucceeded,
				Config:            configFile,
				ResourceType:      resourceDefinition.Properties.ResourceType,
				ProviderType:      providerRegistrationPackage.ProviderType,
			})
			if err != nil {
				fmt.Printf("Failed to insert data: %s", err)
			}
		}()
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
	response.Header().Set(consts.AzureAsyncOperationHeader, getAsyncOperationURI(request.HeaderParameter(consts.RefererHeader), engines.GetAzureAsyncOperationID(request)))
	response.WriteHeader(http.StatusCreated)
	response.Write(responseContent)
}

// DeleteResourceController deletes a resource
func (resourceManager *ResourceManager) DeleteResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	resourcePackage := entities.ResourcePackage{}
	err := resourceManager.ResourceDataProvider.FindPackage(fullyQualifiedResourceID, &resourcePackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			fmt.Sprintf("Resource with id '%s' was not found", fullyQualifiedResourceID))
		return
	}

	if strings.EqualFold(resourcePackage.ProvisioningState, consts.ProvisioningStateAccepted) {
		apierror.WriteErrorToResponse(
			response,
			http.StatusConflict,
			apierror.ClientError,
			apierror.Conflict,
			fmt.Sprintf("Cannot delete Resource with id '%s' as it is being provisioned", fullyQualifiedResourceID))
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
		err := resourceManager.ResourceDataProvider.RemovePackage(fullyQualifiedResourceID)
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

// GetOperationStatusController returns an opeartion status
func (resourceManager *ResourceManager) GetOperationStatusController(request *restful.Request, response *restful.Response) {
	fullyQualifiedOperationStatusID := engines.GetFullyQualifiedOperationStatusID(request)

	// Get Document from collection
	resourcePackage := entities.ResourcePackage{}
	err := resourceManager.ResourceDataProvider.FindPackage(fullyQualifiedOperationStatusID, &resourcePackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			err.Error())
		return
	}

	responseContent, err := json.Marshal(resourcePackage.ToAsyncOperationResult())
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

func getAsyncOperationURI(baseURI string, resourceID string) string {
	return "https://management.azure.com/" + resourceID
}
