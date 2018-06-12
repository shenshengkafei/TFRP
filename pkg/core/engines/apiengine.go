//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package engines

import (
	"TFRP/pkg/core/consts"

	restful "github.com/emicklei/go-restful"
)

const (
	// Subscriptions identifies a customer/an organization
	Subscriptions = "subscriptions"
	// ResourceGroups contains a group of ARM (RP) resources
	ResourceGroups = "resourceGroups"
	// Resources is the name of resource
	Resources = "resources"
	// Providers is the name of provider
	Providers = "providers"
	// ProviderRegistrations is the name of provider registration
	ProviderRegistrations = "providerregistrations"
	// OperationStatus is the operation status
	OperationStatus = "operationstatus"
)

// GetSubscriptionID returns the subscription id if it was on the request else empty string
func GetSubscriptionID(request *restful.Request) string {
	return request.PathParameter(consts.PathSubscriptionIDParameter)
}

// GetAPIVersion returns the api version on the request
func GetAPIVersion(request *restful.Request) string {
	return request.QueryParameter(consts.RequestAPIVersionParameterName)
}

// GetResourceGroupName returns the resourceGroupName if it was on the request else empty string
func GetResourceGroupName(request *restful.Request) string {
	return request.PathParameter(consts.PathResourceGroupNameParameter)
}

// GetResourceName returns the resourceName if it was on the request else empty string
func GetResourceName(request *restful.Request) string {
	return request.PathParameter(consts.PathResourceNameParameter)
}

// GetOperationStatusID returns the operationStatusId if it was on the request else empty string
func GetOperationStatusID(request *restful.Request) string {
	return request.PathParameter(consts.PathOperationStatusParameter)
}

// GetProviderRegistrationName returns the provider registration name if it was on the request else empty string
func GetProviderRegistrationName(request *restful.Request) string {
	return request.PathParameter(consts.PathProviderRegistrationParameter)
}

// GetFullyQualifiedResourceID returns the fully qualified resource id
// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.TerraformOSS/resources/{resource}
func GetFullyQualifiedResourceID(request *restful.Request) string {
	return "/" + Subscriptions + "/" + GetSubscriptionID(request) +
		"/" + ResourceGroups + "/" + GetResourceGroupName(request) +
		"/" + Providers + "/" + consts.TerraformRPNamespace +
		"/" + Resources + "/" + GetResourceName(request)
}

// GetFullyQualifiedProviderRegistrationID returns the fully qualified id of provider registration
// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.TerraformOSS/providerregistrations/{providerRegistration}
func GetFullyQualifiedProviderRegistrationID(request *restful.Request) string {
	return "/" + Subscriptions + "/" + GetSubscriptionID(request) +
		"/" + ResourceGroups + "/" + GetResourceGroupName(request) +
		"/" + Providers + "/" + consts.TerraformRPNamespace +
		"/" + ProviderRegistrations + "/" + GetProviderRegistrationName(request)
}

// GetFullyQualifiedOperationStatusID returns the fully qualified resource operation status id
// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.TerraformOSS/resources/{resource}
func GetFullyQualifiedOperationStatusID(request *restful.Request) string {
	return "/" + Subscriptions + "/" + GetSubscriptionID(request) +
		"/" + ResourceGroups + "/" + GetResourceGroupName(request) +
		"/" + Providers + "/" + consts.TerraformRPNamespace +
		"/" + Resources + "/" + GetOperationStatusID(request)
}

// GetAzureAsyncOperationID returns the operation status id
// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.TerraformOSS/operationstatus/{operationstatusId}
func GetAzureAsyncOperationID(request *restful.Request) string {
	return Subscriptions + "/" + GetSubscriptionID(request) +
		"/" + ResourceGroups + "/" + GetResourceGroupName(request) +
		"/" + Providers + "/" + consts.TerraformRPNamespace +
		"/" + OperationStatus + "/" + GetResourceName(request) +
		"?" + consts.RequestAPIVersionParameterName + "=" + consts.OperationStatusAPIVersion

}
