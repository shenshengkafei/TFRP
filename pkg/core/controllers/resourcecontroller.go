//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package controllers

import (
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/entities"
	"TFRP/pkg/core/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	restful "github.com/emicklei/go-restful"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/mgo.v2/bson"
)

// GetResourceController returns a resource
func GetResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	result := entities.ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := engines.GetProvider(result.ProviderType)

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.Init()
	state.ID = result.StateID

	// Call refresh
	resultState, err := provider.Refresh(info, state)
	if err != nil {
		fmt.Printf("%s", err)
	}

	responseBody, _ := json.Marshal(resultState)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}

// PutResourceController creates/updates a resource
func PutResourceController(request *restful.Request, response *restful.Response) {
	resourceDefinition := entities.ResourceDefinition{}

	rawBody, err := ioutil.ReadAll(request.Request.Body)
	err = json.Unmarshal(rawBody, &resourceDefinition)

	// Get Document from collection
	result := entities.ProviderRegistrationPackage{}
	err = storage.GetProviderRegistrationDataProvider().Find(bson.M{"resourceid": resourceDefinition.Properties.ProviderID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	resourceSpec, _ := json.Marshal(resourceDefinition.Properties.Settings)

	// configFile := getKubernetesTemplateInJson(decoded, resourceDefinition, engines.GetResourceName(request), resourceSpec)
	configFile := getConfigFileInJSON(result.ProviderType, result.Credentials, resourceDefinition, engines.GetResourceName(request), resourceSpec)
	fmt.Printf("%s", configFile)

	provider := engines.GetProvider(result.ProviderType)

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: resourceDefinition.Properties.ResourceType,
	}

	for _, v := range cfg.Resources {
		state := new(terraform.InstanceState)
		state.Init()
		diff, err := provider.Diff(info, state, terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			fmt.Printf("%s", err)
		}

		// Call apply to create resource
		resultState, _ := provider.Apply(info, state, diff)
		fmt.Printf("%s", resultState.ID)

		fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

		// insert Document in collection
		err = storage.GetResourceDataProvider().Insert(&entities.ResourcePackage{
			ResourceID:   fullyQualifiedResourceID,
			StateID:      resultState.ID,
			Config:       configFile,
			ResourceType: resourceDefinition.Properties.ResourceType,
			ProviderType: result.ProviderType,
		})

		if err != nil {
			log.Fatal("Problem inserting data: ", err)
			return
		}
	}

	responseBody, _ := json.Marshal(resourceDefinition)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}

// DeleteResourceController deletes a resource
func DeleteResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	result := entities.ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := engines.GetProvider(result.ProviderType)

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.ID = result.StateID

	diff := new(terraform.InstanceDiff)
	diff.Destroy = true

	// Call apply to delete resource
	resultState, _ := provider.Apply(info, state, diff)

	responseBody, _ := json.Marshal(resultState)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
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
