//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package main

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/controllers"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/storage"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/spf13/pflag"
)

var (
	addr       = pflag.String("insecure-address", ":8080", "The <host>:<port> for insecure (HTTP) serving")
	secureAddr = pflag.String("secure-address", ":443", "The <host>:<port> for secure (HTTPS) serving")
)

func main() {
	pflag.Parse()

	initRoutes()

	go func() {
		go log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Fatal(http.ListenAndServeTLS(":443", "fullchain.pem", "privkey.pem", nil))
}

func initRoutes() {
	// secretEngine := &engines.SecretEngine{
	// 	TenantID:     "72f988bf-86f1-41af-91ab-2d7cd011db47",
	// 	ClientID:     "962c0e07-48bc-4e30-b5fe-b95c11fe486c",
	// 	ClientSecret: "c5c8960e1119c4e0f944",
	// }
	secretEngine := engines.GetSecretEngine()

	storagePassword := secretEngine.GetSecretFromKeyVault(consts.StoragePasswordKVBaseURI, consts.StoragePasswordKVSecretName, consts.StoragePasswordKVSecretVersion)
	providerRegistrationDataProvider := storage.NewProviderRegistrationDataProvider(consts.StorageDatabase, storagePassword)
	resourceDataProvider := storage.NewResourceDataProvider(consts.StorageDatabase, storagePassword)
	providerRegistrationManager := controllers.NewProviderRegistrationManager(providerRegistrationDataProvider)
	resourceManager := controllers.NewResourceManager(providerRegistrationDataProvider, resourceDataProvider)

	webService := new(restful.WebService)
	webService.
		Path(consts.SubscriptionsURLPrefix).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	addSubscriptionOperationRoutes(webService)
	addProvidersOperationRoutes(webService, providerRegistrationManager)
	addResourcesOperationRoutes(webService, resourceManager)

	restful.Add(webService)
}

func addProvidersOperationRoutes(webService *restful.WebService, providerRegistrationManager *controllers.ProviderRegistrationManager) {
	webService.Route(webService.
		GET(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.GetProviderRegistrationController).
		Doc("Get a provider registration").
		Operation(consts.GetProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.PutProviderRegistrationController).
		Doc("Create/update a provider registration").
		Operation(consts.PutProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.DeleteProviderRegistrationController).
		Doc("Delete a provider registration").
		Operation(consts.DeleteProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addResourcesOperationRoutes(webService *restful.WebService, resourceManager *controllers.ResourceManager) {
	webService.Route(webService.
		GET(consts.ResourceOperationRoute).
		To(resourceManager.GetResourceController).
		Doc("Get a resource").
		Operation(consts.GetResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ResourceOperationRoute).
		To(resourceManager.PutResourceController).
		Doc("Create/update a resource").
		Operation(consts.PutResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ResourceOperationRoute).
		To(resourceManager.DeleteResourceController).
		Doc("Delete a resource").
		Operation(consts.DeleteResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addSubscriptionOperationRoutes(webService *restful.WebService) {
	// Subscription operations
	webService.Route(webService.
		GET(consts.SubscriptionResourceOperationRoute).
		To(controllers.GetSubscriptionOperationController).
		Doc("get a subscription").
		Operation(consts.GetSubscriptionControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "identifier of the subscription").DataType("string")))

	webService.Route(webService.
		PUT(consts.SubscriptionResourceOperationRoute).
		To(controllers.PutSubscriptionOperationController).
		Doc("Put or update a subscription").
		Operation(consts.PutSubscriptionControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "identifier of the subscription").DataType("string")))
}
