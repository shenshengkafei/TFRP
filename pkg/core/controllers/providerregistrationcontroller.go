//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

// This file is for code about storing and retrieving api tracking
// info from a context struct

package controllers

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/entities"
	"TFRP/pkg/core/storage"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	restful "github.com/emicklei/go-restful"
)

// GetProviderRegistrationController returns a provider registration
func GetProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err := storage.GetProviderRegistrationDataProvider().Find(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	responseContent, err := json.Marshal(providerRegistrationPackage)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// PutProviderRegistrationController create a new provider registration
func PutProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	providerRegistrationDefinition := entities.ProviderRegistrationDefinition{}
	rawBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}
	err = json.Unmarshal(rawBody, &providerRegistrationDefinition)
	credentials, err := json.Marshal(providerRegistrationDefinition.Properties.Settings)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	if strings.EqualFold(consts.KubernetesProvider, providerRegistrationDefinition.Properties.ProviderType) {
		credentials = getKubernetesProviderCredentials(providerRegistrationDefinition.Properties.Settings)
	}

	fmt.Printf("%s", string(credentials))
	// insert Document in collection
	err = storage.GetProviderRegistrationDataProvider().Insert(&entities.ProviderRegistrationPackage{
		ResourceID:   fullyQualifiedResourceID,
		ProviderType: providerRegistrationDefinition.Properties.ProviderType,
		Credentials:  credentials,
	})

	if err != nil {
		log.Fatal("Problem inserting data: ", err)
		return
	}

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err = storage.GetProviderRegistrationDataProvider().Find(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	responseContent, _ := json.Marshal(providerRegistrationPackage)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

// DeleteProviderRegistrationController removes a provider registration
func DeleteProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	// Get Document from collection
	providerRegistrationPackage := entities.ProviderRegistrationPackage{}
	err := storage.GetProviderRegistrationDataProvider().Find(fullyQualifiedResourceID, &providerRegistrationPackage)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	err = storage.GetProviderRegistrationDataProvider().Remove(fullyQualifiedResourceID)
	if err != nil {
		log.Fatal("Error deleting record: ", err)
		return
	}

	responseContent, _ := json.Marshal(providerRegistrationPackage)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseContent)
}

func getKubernetesProviderCredentials(credentials interface{}) []byte {
	kubeCredentials := &entities.KubernetesProviderCredential{}

	byteData, err := json.Marshal(credentials)
	if err != nil {
		fmt.Printf("%s", err)
	}

	err = json.Unmarshal(byteData, &kubeCredentials)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", kubeCredentials.InlineConfig)

	decodedConfig, err := base64.StdEncoding.DecodeString(kubeCredentials.InlineConfig)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", decodedConfig)

	decodedKubeCredentials := &entities.KubernetesProviderCredential{
		InlineConfig: string(decodedConfig),
	}

	result, err := json.Marshal(decodedKubeCredentials)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return result
}
