//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

// This file is for code about storing and retrieving api tracking
// info from a context struct

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
	"strings"

	restful "github.com/emicklei/go-restful"
)

// ProviderRegistrationManager is the provider registeration manager
type ProviderRegistrationManager struct {
	BaseHandler
}

// NewProviderRegistrationManager create a new provider registration manager
func NewProviderRegistrationManager(providerRegistrationDataProvider *storage.ProviderRegistrationDataProvider) (providerRegistrationManager *ProviderRegistrationManager) {
	providerRegistrationManager = new(ProviderRegistrationManager)
	providerRegistrationManager.ProviderRegistrationDataProvider = providerRegistrationDataProvider
	return providerRegistrationManager
}

// GetProviderRegistrationController returns a provider registration
func (providerRegistrationManager *ProviderRegistrationManager) GetProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err := providerRegistrationManager.ProviderRegistrationDataProvider.FindPackage(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			err.Error())
		return
	}

	responseContent, err := json.Marshal(providerRegistrationPackage.ToDefinition())
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize provider registration package: %s", err.Error()))
		return
	}
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// PutProviderRegistrationController create a new provider registration
func (providerRegistrationManager *ProviderRegistrationManager) PutProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	providerRegistrationDefinition := entities.ProviderRegistrationDefinition{}
	rawBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("The request body is invalid: %s", err.Error()))
		return
	}

	err = json.Unmarshal(rawBody, &providerRegistrationDefinition)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusBadRequest,
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Failed to deserialize request content: %s", err.Error()))
		return
	}

	settings, err := json.Marshal(providerRegistrationDefinition.Properties.Settings)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize provider registration settings: %s", err.Error()))
		return
	}

	apierr := engines.ValidateProviderRegistrationDefinition(&providerRegistrationDefinition)
	if apierr != nil {
		apierror.WriteErrorToResponseWitAPIError(
			response,
			http.StatusBadRequest,
			apierr)
		return
	}

	fmt.Printf("%s", string(settings))
	// insert Document in collection
	err = providerRegistrationManager.ProviderRegistrationDataProvider.InsertPackage(&entities.ProviderRegistrationPackage{
		ResourceID:   fullyQualifiedResourceID,
		ProviderType: strings.ToLower(providerRegistrationDefinition.Properties.ProviderType),
		Settings:     settings,
	})

	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to insert data: %s", err.Error()))
		return
	}

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err = providerRegistrationManager.ProviderRegistrationDataProvider.FindPackage(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to find data: %s", err.Error()))
		return
	}

	responseContent, err := json.Marshal(providerRegistrationPackage.ToDefinition())
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize response content: %s", err.Error()))
		return
	}

	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// DeleteProviderRegistrationController removes a provider registration
func (providerRegistrationManager *ProviderRegistrationManager) DeleteProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err := providerRegistrationManager.ProviderRegistrationDataProvider.FindPackage(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			fmt.Sprintf("Provider '%s' was not found.", fullyQualifiedResourceID))
		return
	}

	err = providerRegistrationManager.ProviderRegistrationDataProvider.RemovePackage(fullyQualifiedResourceID)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to delete data: %s", err.Error()))
		return
	}

	response.WriteHeader(http.StatusOK)
}

// PostProviderRegistrationListSettings returns settings of a provider registration
func (providerRegistrationManager *ProviderRegistrationManager) PostProviderRegistrationListSettings(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err := providerRegistrationManager.ProviderRegistrationDataProvider.FindPackage(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusNotFound,
			apierror.ClientError,
			apierror.NotFound,
			err.Error())
		return
	}

	responseContent, err := json.Marshal(providerRegistrationPackage.ToListSettingsDefinition())
	if err != nil {
		apierror.WriteErrorToResponse(
			response,
			http.StatusInternalServerError,
			apierror.InternalError,
			apierror.InternalOperationError,
			fmt.Sprintf("Failed to serialize provider registration package: %s", err.Error()))
		return
	}
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}
