//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package main

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/controllers"
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
	// router := mux.NewRouter()
	// router.NotFoundHandler = http.HandlerFunc(NotFound)
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/provider/{provider}", putProvider).Methods("PUT")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", getResource).Methods("GET")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", putResource).Methods("PUT")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", deleteResource).Methods("DELETE")
	// //log.Fatal(http.ListenAndServeTLS(":443", "fullchain.pem", "privkey.pem", router))
	// log.Fatal(http.ListenAndServe(":8080", router))

	pflag.Parse()

	initRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initRoutes() {
	webService := new(restful.WebService)
	webService.
		Path(consts.SubscriptionsURLPrefix).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	addProvidersOperationRoutes(webService)
	addResourcesOperationRoutes(webService)

	restful.Add(webService)
}

func addProvidersOperationRoutes(webService *restful.WebService) {
	webService.Route(webService.
		GET(consts.ProviderRegistrationOperationRoute).
		To(controllers.GetProviderRegistrationController).
		Doc("Get a provider registration").
		Operation(consts.GetProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ProviderRegistrationOperationRoute).
		To(controllers.PutProviderRegistrationController).
		Doc("Create/update a provider registration").
		Operation(consts.PutProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ProviderRegistrationOperationRoute).
		To(controllers.DeleteProviderRegistrationController).
		Doc("Delete a provider registration").
		Operation(consts.DeleteProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addResourcesOperationRoutes(webService *restful.WebService) {
	webService.Route(webService.
		GET(consts.ResourceOperationRoute).
		To(controllers.GetResourceController).
		Doc("Get a resource").
		Operation(consts.GetResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ResourceOperationRoute).
		To(controllers.PutResourceController).
		Doc("Create/update a resource").
		Operation(consts.PutResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ResourceOperationRoute).
		To(controllers.DeleteResourceController).
		Doc("Delete a resource").
		Operation(consts.DeleteResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}
