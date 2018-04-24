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
	"gopkg.in/mgo.v2/bson"
)

// PutProviderRegistrationController create a new provider registration
func PutProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	providerRegistrationDefinition := entities.ProviderRegistrationDefinition{}
	rawBody, err := ioutil.ReadAll(request.Request.Body)
	err = json.Unmarshal(rawBody, &providerRegistrationDefinition)
	credentials, _ := json.Marshal(providerRegistrationDefinition.Properties.Settings)

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
	result := entities.ProviderRegistrationPackage{}
	err = storage.GetProviderRegistrationDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	responseBody, _ := json.Marshal(result)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
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
